package mqtt

import (
	"math"

	"github.com/eunnseo/AirPost/application/domain/model"
)

// earthRadiusMeters is the mean Earth radius used for the local NE projection.
const earthRadiusMeters = 6371000.0

// localNorthEast projects the delta from an origin (lat/lon in degrees) to a
// target point onto a local tangent plane, returning north/east meters. This is
// the equirectangular approximation, which is accurate for the short station-
// to-tag distances AirPost flies and keeps the function dependency-free/pure.
func localNorthEast(originLat, originLon, targetLat, targetLon float64) (north, east float64) {
	rad := math.Pi / 180
	north = (targetLat - originLat) * rad * earthRadiusMeters
	east = (targetLon - originLon) * rad * earthRadiusMeters * math.Cos(originLat*rad)
	return north, east
}

// BuildDeliveryRequest builds the MQTT request JSON payload for a delivery.
// takeoff is the station the drone lifts off from, pickup is the station holding
// the parcel (== takeoff unless the drone is ferried in from elsewhere), landing
// is the chosen nearest destination station, and tag is the recipient drop point.
// The delivery point is encoded as north/east meters relative to the PICKUP
// station, since the parcel always travels from where it is picked up. cruise is
// the assigned cruise altitude in meters. This is a pure function so the mapping
// is unit-testable without a broker or database.
func BuildDeliveryRequest(d *model.Delivery, takeoff, pickup, landing, tag *model.Node, cruise float64) DeliveryRequest {
	north, east := localNorthEast(pickup.LocLat, pickup.LocLon, tag.LocLat, tag.LocLon)
	return DeliveryRequest{
		OrderID:   d.OrderNum,
		DroneID:   d.DroneID,
		TakeoffID: takeoff.ID,
		PickupID:  pickup.ID,
		DeliverN:  north,
		DeliverE:  east,
		LandingID: landing.ID,
		Cruise:    cruise,
	}
}
