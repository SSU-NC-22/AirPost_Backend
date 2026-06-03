// AirPost
package main

import (
	"log"
	"os"

	"github.com/eunnseo/AirPost/application/dataService/sql"
	"github.com/eunnseo/AirPost/application/delivery"
	deliverymqtt "github.com/eunnseo/AirPost/application/delivery/mqtt"
	"github.com/eunnseo/AirPost/application/docs"
	"github.com/eunnseo/AirPost/application/domain/model"
	"github.com/eunnseo/AirPost/application/domain/repository"
	"github.com/eunnseo/AirPost/application/rest/handler"
	"github.com/eunnseo/AirPost/application/setting"
	"github.com/eunnseo/AirPost/application/usecase"
	"github.com/eunnseo/AirPost/application/usecase/eventUsecase"
	"github.com/eunnseo/AirPost/application/usecase/registUsecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

func main() {
	sql.Setup()

	sir := sql.NewSinkRepo()
	ndr := sql.NewNodeRepo()
	lgr := sql.NewLogicRepo()
	lsr := sql.NewLogicServiceRepo()
	tpr := sql.NewTopicRepo()

	dlr := sql.NewDeliveryRepo()
	ptr := sql.NewPathRepo()
	sdr := sql.NewStationDroneRepo()

	ru := registUsecase.NewRegistUsecase(sir, ndr, lgr, lsr, tpr, dlr, ptr, sdr)
	eu := eventUsecase.NewEventUsecase(sir, lsr)

	h := handler.NewHandler(ru, eu)

	r := gin.Default()
	// CORS: a wildcard origin "*" combined with AllowCredentials:true is
	// invalid per the Fetch spec and is rejected by browsers. Restrict to a
	// real UI origin (configurable via the UI_ORIGIN env var).
	uiOrigin := os.Getenv("UI_ORIGIN")
	if uiOrigin == "" {
		uiOrigin = "http://localhost:3000"
	}
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{uiOrigin}
	config.AllowCredentials = true
	// The default allow-headers omit Authorization, so the browser's preflight blocked every
	// authenticated admin request (the JWT rides in the Authorization header). Allow it.
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// swagger
	docs.SwaggerInfo.Title = "AirPost application API"
	docs.SwaggerInfo.Description = "This is a registration server for AirPost UI."
	docs.SwaggerInfo.Version = "0.1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Wire MQTT so RegistDelivery publishes flight requests and status updates
	// (delivered -> email) are handled. A broker failure is logged but does not
	// stop the API: deliveries are still recorded.
	if dispatcher := setupDelivery(ru); dispatcher != nil {
		h.SetDeliveryDispatcher(dispatcher)
	}

	// Public route: login issues JWTs and must not itself require one.
	r.POST("/auth/login", h.Login)

	// /event is service-to-service infrastructure (logic-core self-registers here at startup with
	// no user JWT). Register it BEFORE the global JWT middleware so it is not gated by user auth —
	// otherwise logic-core gets 401 and crash-loops. Reached only over the internal network in prod.
	setEventRoute(r, h)

	// All routes below require a valid JWT when auth is enabled (default ON).
	// Set AUTH_ENABLED=0 to disable for tests/dev.
	if handler.AuthEnabled() {
		r.Use(handler.JWTAuthMiddleware())
	}

	setRegistrationRoute(r, h)
	initTopic(tpr)

	initDroneSink(sir, eu)
	initStationSink(sir, eu)
	initTagSink(sir, eu)

	// Seed a usable demo topology (stations, drones, tags, station-drone links,
	// paths) once the sinks exist, so a fresh `compose up` flies real sorties with
	// no manual setup. Idempotent; skipped when SEED_DEMO=0.
	sql.Seed()

	log.Fatal(r.Run(setting.Appsetting.Server))
}

// adminOnly gates a route group to admin callers when auth is enabled. When
// auth is off (tests/dev) it is a no-op so the group stays usable.
func adminOnly() gin.HandlerFunc {
	if handler.AuthEnabled() {
		return handler.RequireRole(handler.RoleAdmin)
	}
	return func(c *gin.Context) { c.Next() }
}

func setEventRoute(r *gin.Engine, h *handler.Handler) {
	// Logic-service registration is service-to-service infrastructure (logic-core calls this at
	// startup with no user JWT). It must NOT sit behind the user adminOnly() gate — doing so
	// returned 401 and crash-looped logic-core. In a real deployment this route is reached only
	// over the internal network, not exposed publicly.
	event := r.Group("/event")
	{
		event.POST("", h.RegistLogicService)
	}
}

func setRegistrationRoute(r *gin.Engine, h *handler.Handler) {
	regist := r.Group("/regist")
	{
		// Infrastructure CRUD (sinks, nodes, logic, topics) is admin only.
		sink := regist.Group("/sink", adminOnly())
		{
			sink.GET("", h.ListSinks)
			sink.POST("", h.RegistSink)
			sink.DELETE("/:id", h.UnregistSink)
		}
		node := regist.Group("/node", adminOnly())
		{
			node.GET("", h.ListNodes)
			node.GET("/:sinkid", h.ListNodesBySink)
			node.POST("", h.RegistNode)
			node.POST("/update", h.UpdateNodeLoc)
			node.DELETE("/:id", h.UnregistNode)
		}
		logic := regist.Group("/logic", adminOnly())
		{
			logic.GET("", h.ListLogics)
			logic.POST("", h.RegistLogic) // << 프론트에서
			logic.DELETE("/:id", h.UnregistLogic)
		}
		logicService := regist.Group("/logic-service", adminOnly())
		{
			logicService.GET("", h.ListLogicServices)
			logicService.DELETE("/:id", h.UnregistLogicService)
		}
		topic := regist.Group("/topic", adminOnly())
		{
			topic.GET("", h.ListTopics)
			topic.POST("", h.RegistTopic)
			topic.DELETE("/:id", h.UnregistTopic)
		}
		// Deliveries and tracking are usable by any authenticated user; the
		// handlers enforce per-record ownership (or admin) on reads.
		delivery := regist.Group("/delivery")
		{
			delivery.GET("/:orderNum", h.GetDroneID)
			delivery.POST("", h.RegistDelivery)
		}
		tracking := regist.Group("/tracking")
		{
			tracking.GET("/:orderNum", h.GetTracking)
		}
	}
}

// setupDelivery connects to the MQTT broker, subscribes to delivery status
// updates, and returns a Dispatcher for publishing flight requests. It returns
// nil (logging the cause) if the broker is unreachable, so the API still runs.
func setupDelivery(ru usecase.RegistUsecase) *delivery.Dispatcher {
	client, err := deliverymqtt.NewClient("airpost-application")
	if err != nil {
		log.Printf("delivery: MQTT disabled, broker unavailable: %v", err)
		return nil
	}

	dispatcher := delivery.NewDispatcher(client, ru, sql.NewBusyRepo())
	if err := client.SubscribeStatus(dispatcher.HandleStatus); err != nil {
		log.Printf("delivery: status subscription failed: %v", err)
	}
	return dispatcher
}

func initTopic(tpr repository.TopicRepo) {
	if setting.Topicsetting.Name != "" {
		t := model.Topic{
			Name:         setting.Topicsetting.Name,
			Partitions:   setting.Topicsetting.Partitions,
			Replications: setting.Topicsetting.Replications,
		}
		tpr.Create(&t)
	}
}

func initDroneSink(sir repository.SinkRepo, eu usecase.EventUsecase) {
	if setting.DroneSinksetting.Name != "" {
		s := model.Sink{
			Name:		setting.DroneSinksetting.Name,
			Addr:		setting.DroneSinksetting.Addr,
			TopicID:	setting.DroneSinksetting.TopicID,
		}
		sir.Create(&s)
		eu.CreateSinkEvent(&s)
	}
}

func initStationSink(sir repository.SinkRepo, eu usecase.EventUsecase) {
	if setting.StationSinksetting.Name != "" {
		s := model.Sink{
			Name:		setting.StationSinksetting.Name,
			Addr:		setting.StationSinksetting.Addr,
			TopicID:	setting.StationSinksetting.TopicID,
		}
		sir.Create(&s)
		eu.CreateSinkEvent(&s)
	}
}

func initTagSink(sir repository.SinkRepo, eu usecase.EventUsecase) {
	if setting.TagSinksetting.Name != "" {
		s := model.Sink{
			Name:		setting.TagSinksetting.Name,
			Addr:		setting.TagSinksetting.Addr,
			TopicID:	setting.TagSinksetting.TopicID,
		}
		sir.Create(&s)
		eu.CreateSinkEvent(&s)
	}
}

