// Package delivery wires the delivery flow to the MQTT flight service: it
// publishes flight requests and reacts to status updates (notably sending the
// "delivered" email). It depends only on small interfaces so it stays testable.
package delivery

import (
	"log"
	"os"
	"strconv"
	"sync"

	deliverymqtt "github.com/eunnseo/AirPost/application/delivery/mqtt"
	"github.com/eunnseo/AirPost/application/delivery/notify"
	"github.com/eunnseo/AirPost/application/domain/model"
)

// Altitude-band deconfliction: each concurrent mission cruises in its own band so
// two airborne drones never share an altitude. baseCruise is the lowest band and
// bandGap the vertical separation between bands (meters). CRUISE_ALTITUDE overrides
// the base. Takeoff/landing happen at distinct stations, so only the cruise leg
// needs separating.
const (
	defaultCruiseAltitude = 30.0
	bandGap               = 6.0
)

// requestPublisher publishes flight requests; satisfied by mqtt.Client.
type requestPublisher interface {
	PublishRequest(deliverymqtt.DeliveryRequest) error
}

// deliveryLookup resolves a delivery by its order number for status handling.
type deliveryLookup interface {
	GetDeliveryByOrderNum(orderNum string) (model.Delivery, error)
}

// Dispatcher publishes delivery requests and handles inbound status updates. It
// also acts as the "control tower": it tracks each in-flight mission's altitude
// band and frees it on completion, keeping concurrent drones vertically separated.
type Dispatcher struct {
	publisher requestPublisher
	lookup    deliveryLookup
	smtp      notify.SMTPConfig

	mu         sync.Mutex
	bands      map[string]int // order_id -> altitude band slot
	used       []bool         // band slot -> in use
	orderDrone map[string]int // order_id -> drone id (to free the drone on landing)
	busy       map[int]bool   // drone id -> in flight (skipped when assigning)
}

// NewDispatcher builds a Dispatcher from a publisher and a delivery lookup.
func NewDispatcher(publisher requestPublisher, lookup deliveryLookup) *Dispatcher {
	return &Dispatcher{
		publisher:  publisher,
		lookup:     lookup,
		smtp:       notify.LoadSMTPConfig(),
		bands:      make(map[string]int),
		orderDrone: make(map[string]int),
		busy:       make(map[int]bool),
	}
}

// Dispatch maps the delivery and its resolved stations/tag to the MQTT request,
// assigns a free altitude band for collision deconfliction, and publishes it so
// the flight service starts the mission. takeoff is where the drone lifts off,
// pickup the parcel's source station (may differ from takeoff when ferrying),
// landing the nearest station to the drop, and tag the drop point.
func (d *Dispatcher) Dispatch(delivery *model.Delivery, takeoff, pickup, landing, tag *model.Node) error {
	cruise := d.reserveBand(delivery.OrderNum)
	d.markBusy(delivery.OrderNum, delivery.DroneID)
	req := deliverymqtt.BuildDeliveryRequest(delivery, takeoff, pickup, landing, tag, cruise)
	if err := d.publisher.PublishRequest(req); err != nil {
		d.releaseBand(delivery.OrderNum) // don't leak the band/drone if publish failed
		d.freeDrone(delivery.OrderNum)
		return err
	}
	return nil
}

// markBusy records the drone flying an order so it is not assigned to another
// order until it lands.
func (d *Dispatcher) markBusy(orderID string, droneID int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.orderDrone[orderID] = droneID
	d.busy[droneID] = true
}

// freeDrone releases the drone held by an order so the fleet can reuse it.
func (d *Dispatcher) freeDrone(orderID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if id, ok := d.orderDrone[orderID]; ok {
		delete(d.busy, id)
		delete(d.orderDrone, orderID)
	}
}

// IsDroneBusy reports whether a drone is currently flying a mission. A nil
// dispatcher (no broker in dev/tests) reports everything idle.
func (d *Dispatcher) IsDroneBusy(droneID int) bool {
	if d == nil {
		return false
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.busy[droneID]
}

// reserveBand assigns the lowest free altitude band to an order and returns its
// cruise altitude in meters.
func (d *Dispatcher) reserveBand(orderID string) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	if slot, ok := d.bands[orderID]; ok { // idempotent on re-dispatch
		return bandAltitude(slot)
	}
	slot := 0
	for slot < len(d.used) && d.used[slot] {
		slot++
	}
	if slot == len(d.used) {
		d.used = append(d.used, true)
	} else {
		d.used[slot] = true
	}
	d.bands[orderID] = slot
	return bandAltitude(slot)
}

// releaseBand frees the altitude band held by an order, if any.
func (d *Dispatcher) releaseBand(orderID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if slot, ok := d.bands[orderID]; ok {
		if slot < len(d.used) {
			d.used[slot] = false
		}
		delete(d.bands, orderID)
	}
}

// bandAltitude maps a band slot to its cruise altitude.
func bandAltitude(slot int) float64 {
	return cruiseAltitude() + float64(slot)*bandGap
}

// HandleStatus reacts to a status update from the flight service. On a
// successful delivery it emails the recipient. Other states are logged for now;
// richer record updates can hook in here without touching the transport.
func (d *Dispatcher) HandleStatus(status deliverymqtt.DeliveryStatus) {
	log.Printf("delivery %s state=%s result=%s", status.OrderID, status.State, status.Result)

	// Free the mission's altitude band and drone once it lands (terminal result),
	// so both are available to the next mission. Mid-flight states keep them held.
	if status.Result != "" {
		d.releaseBand(status.OrderID)
		d.freeDrone(status.OrderID)
	}

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
