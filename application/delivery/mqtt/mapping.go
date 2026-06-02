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
// takeoff is the source station, landing is the chosen nearest destination
// station, and tag is the recipient drop point; the delivery point is encoded
// as north/east meters relative to the takeoff station. cruise is the cruise
// altitude in meters. This is a pure function so the mapping is unit-testable
// without a broker or database.
func BuildDeliveryRequest(d *model.Delivery, takeoff, landing, tag *model.Node, cruise float64) DeliveryRequest {
	north, east := localNorthEast(takeoff.LocLat, takeoff.LocLon, tag.LocLat, tag.LocLon)
	return DeliveryRequest{
		OrderID:   d.OrderNum,
		TakeoffID: takeoff.ID,
		DeliverN:  north,
		DeliverE:  east,
		LandingID: landing.ID,
		Cruise:    cruise,
	}
}
