// Package delivery wires the delivery flow to the MQTT flight service: it
// publishes flight requests and reacts to status updates (notably sending the
// "delivered" email). It depends only on small interfaces so it stays testable.
package delivery

import (
	"log"
	"os"
	"strconv"

	deliverymqtt "github.com/eunnseo/AirPost/application/delivery/mqtt"
	"github.com/eunnseo/AirPost/application/delivery/notify"
	"github.com/eunnseo/AirPost/application/domain/model"
)

// defaultCruiseAltitude is used when CRUISE_ALTITUDE is unset (meters).
const defaultCruiseAltitude = 30.0

// requestPublisher publishes flight requests; satisfied by mqtt.Client.
type requestPublisher interface {
	PublishRequest(deliverymqtt.DeliveryRequest) error
}

// deliveryLookup resolves a delivery by its order number for status handling.
type deliveryLookup interface {
	GetDeliveryByOrderNum(orderNum string) (model.Delivery, error)
}

// Dispatcher publishes delivery requests and handles inbound status updates.
type Dispatcher struct {
	publisher requestPublisher
	lookup    deliveryLookup
	smtp      notify.SMTPConfig
}

// NewDispatcher builds a Dispatcher from a publisher and a delivery lookup.
func NewDispatcher(publisher requestPublisher, lookup deliveryLookup) *Dispatcher {
	return &Dispatcher{
		publisher: publisher,
		lookup:    lookup,
		smtp:      notify.LoadSMTPConfig(),
	}
}

// Dispatch maps the delivery and its resolved stations/tag to the MQTT request
// and publishes it, so the flight service starts the mission.
func (d *Dispatcher) Dispatch(delivery *model.Delivery, takeoff, landing, tag *model.Node) error {
	req := deliverymqtt.BuildDeliveryRequest(delivery, takeoff, landing, tag, cruiseAltitude())
	return d.publisher.PublishRequest(req)
}

// HandleStatus reacts to a status update from the flight service. On a
// successful delivery it emails the recipient. Other states are logged for now;
// richer record updates can hook in here without touching the transport.
func (d *Dispatcher) HandleStatus(status deliverymqtt.DeliveryStatus) {
	log.Printf("delivery %s state=%s result=%s", status.OrderID, status.State, status.Result)
	if !status.IsDelivered() {
		return
	}

	delivery, err := d.lookup.GetDeliveryByOrderNum(status.OrderID)
	if err != nil {
		log.Printf("delivery %s: cannot load record for notification: %v", status.OrderID, err)
		return
	}
	if err := notify.SendDeliveredEmail(d.smtp, delivery.Email, delivery.OrderNum); err != nil {
		log.Printf("delivery %s: delivered email failed: %v", status.OrderID, err)
		return
	}
	log.Printf("delivery %s: delivered email sent to %s", status.OrderID, delivery.Email)
}

// cruiseAltitude resolves the cruise altitude (meters) from the environment.
func cruiseAltitude() float64 {
	if v := os.Getenv("CRUISE_ALTITUDE"); v != "" {
		if alt, err := strconv.ParseFloat(v, 64); err == nil {
			return alt
		}
	}
	return defaultCruiseAltitude
}
