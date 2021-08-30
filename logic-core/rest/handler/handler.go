package handler

import (
	"log"
	"net/http"

	"github.com/eunnseo/AirPost/logic-core/adapter"
	"github.com/eunnseo/AirPost/logic-core/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	evuc usecase.EventUsecase
	lcuc usecase.LogicCoreUsecase
}

func NewHandler(evuc usecase.EventUsecase, lcuc usecase.LogicCoreUsecase) *Handler {
	return &Handler{
		evuc: evuc,
		lcuc: lcuc,
	}
}

/**************************************************************/
/* Sink handler                                               */
/**************************************************************/
func (h *Handler) CreateSink(c *gin.Context) {
	var addr adapter.SinkAddr
	log.Println("in CreateSink")
	if err := c.ShouldBindJSON(&addr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.lcuc.AppendSinkAddr(&addr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, addr)

	// var an adapter.Node
	// if err := c.ShouldBindJSON(&an); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// if err := h.evuc.CreateNode(&an, an.Sink.Name); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, an)
}

func (h *Handler) DeleteSink(c *gin.Context) {
	var nl []adapter.Node
	if err := c.ShouldBindJSON(&nl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.evuc.DeleteSink(nl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nl)
}

/**************************************************************/
/* Node handler                                               */
/**************************************************************/
func (h *Handler) CreateNode(c *gin.Context) {
	var an adapter.Node
	if err := c.ShouldBindJSON(&an); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.evuc.CreateNode(&an, an.Sink.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, an)
}

func (h *Handler) DeleteNode(c *gin.Context) {
	var an adapter.Node
	if err := c.ShouldBindJSON(&an); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.evuc.DeleteNode(&an); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, an)
}

/**************************************************************/
/* Logic handler                                              */
/**************************************************************/
func (h *Handler) CreateLogic(c *gin.Context) {
	var al adapter.Logic

	if err := c.ShouldBindJSON(&al); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("in createLogic, logic = ", al)
	if err := h.evuc.CreateLogic(&al); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, al)
}

func (h *Handler) DeleteLogic(c *gin.Context) {
	var al adapter.Logic
	if err := c.ShouldBindJSON(&al); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.evuc.DeleteLogic(&al); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, al)
}

/**************************************************************/
/* Delivery handler                                           */
/**************************************************************/
func (h *Handler) CreateDelivery(c *gin.Context) {
	var ad adapter.Delivery

	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("in CreateDelivery, ad = ", ad)
	if err := h.evuc.CreateDelivery(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ad)
}
