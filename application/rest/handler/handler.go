package handler

import (
	"github.com/eunnseo/AirPost/application/delivery"
	"github.com/eunnseo/AirPost/application/domain/model"
	"github.com/eunnseo/AirPost/application/usecase"
)

// deliveryDispatcher publishes a delivery's flight request over MQTT. It is an
// interface so handlers stay testable and the broker is optional in dev/tests.
type deliveryDispatcher interface {
	Dispatch(d *model.Delivery, takeoff, landing, tag *model.Node) error
}

type Handler struct {
	ru usecase.RegistUsecase
	eu usecase.EventUsecase
	// dispatcher is nil when MQTT is not configured; delivery still succeeds.
	dispatcher deliveryDispatcher
}

func NewHandler(ru usecase.RegistUsecase, eu usecase.EventUsecase) *Handler {
	return &Handler{
		ru: ru,
		eu: eu,
	}
}

// SetDeliveryDispatcher wires the MQTT delivery dispatcher after construction so
// broker setup (which may fail/be absent in dev) does not block handler wiring.
func (h *Handler) SetDeliveryDispatcher(d *delivery.Dispatcher) {
	h.dispatcher = d
}
