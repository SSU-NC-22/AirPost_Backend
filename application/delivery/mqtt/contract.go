// Package mqtt carries the shared delivery messaging contract between the
// AirPost backend and the flight/simulation service. The JSON shapes here MUST
// stay byte-compatible with simulation/tests/airpost_service.py and the topic
// names below; both sides are wired to these exact strings.
package mqtt

// Topics on which delivery requests are published and statuses are received.
const (
	RequestTopic = "airpost/delivery/request"
	StatusTopic  = "airpost/delivery/status"
)

// DeliveryRequest is published to RequestTopic to start a flight. deliver_N /
// deliver_E are local north/east offsets in meters from the takeoff station;
// takeoff_id / landing_id are station node IDs; cruise is the cruise altitude.
type DeliveryRequest struct {
	OrderID   string  `json:"order_id"`
	TakeoffID int     `json:"takeoff_id"`
	DeliverN  float64 `json:"deliver_N"`
	DeliverE  float64 `json:"deliver_E"`
	LandingID int     `json:"landing_id"`
	Cruise    float64 `json:"cruise"`
}

// DeliveryStatus is received on StatusTopic as the flight progresses. result is
// "PASS"/"FAIL" only on terminal states; deliver_err / land_err are meters.
type DeliveryStatus struct {
	OrderID   string  `json:"order_id"`
	State     string  `json:"state"`
	DeliverErr float64 `json:"deliver_err"`
	LandErr   float64 `json:"land_err"`
	Result    string  `json:"result"`
}

// IsDelivered reports whether the status marks a successful delivery completion,
// i.e. the point at which the recipient should be notified by email. The sim
// publishes state=="delivered" mid-flight and a terminal state with
// result=="PASS"; either signal counts as a successful delivery.
func (s DeliveryStatus) IsDelivered() bool {
	return s.State == "delivered" || s.Result == "PASS"
}
