package mqtt

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/eunnseo/AirPost/application/domain/model"
)

func TestBuildDeliveryRequestMapping(t *testing.T) {
	d := &model.Delivery{OrderNum: "AP-123", DroneID: 51}
	takeoff := &model.Node{ID: 2, LocLat: 37.0, LocLon: 127.0} // drone sits here (== pickup)
	landing := &model.Node{ID: 7, LocLat: 37.5, LocLon: 127.5} // nearest dest station
	tag := &model.Node{ID: 9, LocLat: 37.001, LocLon: 127.0}   // drop point, ~111m north

	req := BuildDeliveryRequest(d, takeoff, takeoff, landing, tag, 30)

	if req.OrderID != "AP-123" {
		t.Errorf("OrderID = %q, want AP-123", req.OrderID)
	}
	if req.TakeoffID != 2 {
		t.Errorf("TakeoffID = %d, want 2", req.TakeoffID)
	}
	if req.LandingID != 7 {
		t.Errorf("LandingID = %d, want 7", req.LandingID)
	}
	if req.Cruise != 30 {
		t.Errorf("Cruise = %v, want 30", req.Cruise)
	}
	// 0.001 deg latitude north of the takeoff origin is ~111.2 m north, ~0 east.
	if math.Abs(req.DeliverN-111.2) > 1.0 {
		t.Errorf("DeliverN = %.2f m, want ~111.2 m", req.DeliverN)
	}
	if math.Abs(req.DeliverE) > 0.5 {
		t.Errorf("DeliverE = %.2f m, want ~0 m", req.DeliverE)
	}
}

// TestBuildDeliveryRequestEastward verifies the east projection scales by the
// cosine of latitude (a degree of longitude is shorter away from the equator).
func TestBuildDeliveryRequestEastward(t *testing.T) {
	d := &model.Delivery{OrderNum: "AP-east"}
	takeoff := &model.Node{ID: 1, LocLat: 60.0, LocLon: 10.0}
	tag := &model.Node{ID: 2, LocLat: 60.0, LocLon: 10.001} // 0.001 deg east

	req := BuildDeliveryRequest(d, takeoff, takeoff, takeoff, tag, 25)

	if math.Abs(req.DeliverN) > 0.5 {
		t.Errorf("DeliverN = %.2f m, want ~0 m", req.DeliverN)
	}
	// 0.001 deg lon at 60N ~= 111.2 * cos(60) ~= 55.6 m east.
	if math.Abs(req.DeliverE-55.6) > 1.0 {
		t.Errorf("DeliverE = %.2f m, want ~55.6 m", req.DeliverE)
	}
}

// TestBuildDeliveryRequestFerry verifies that when the drone is ferried in from a
// different station, the delivery offset is measured from the PICKUP station (not
// the take-off station) and both station ids are carried on the request.
func TestBuildDeliveryRequestFerry(t *testing.T) {
	d := &model.Delivery{OrderNum: "AP-ferry", DroneID: 53}
	takeoff := &model.Node{ID: 3, LocLat: 37.2, LocLon: 127.3} // drone ferried from here
	pickup := &model.Node{ID: 2, LocLat: 37.0, LocLon: 127.0}  // parcel source
	landing := &model.Node{ID: 7, LocLat: 37.5, LocLon: 127.5}
	tag := &model.Node{ID: 9, LocLat: 37.001, LocLon: 127.0} // ~111 m north of pickup

	req := BuildDeliveryRequest(d, takeoff, pickup, landing, tag, 42)

	if req.DroneID != 53 {
		t.Errorf("DroneID = %d, want 53", req.DroneID)
	}
	if req.TakeoffID != 3 || req.PickupID != 2 {
		t.Errorf("Takeoff/Pickup = %d/%d, want 3/2", req.TakeoffID, req.PickupID)
	}
	// Offset is from the pickup station, so ~111 m north / ~0 east regardless of takeoff.
	if math.Abs(req.DeliverN-111.2) > 1.0 || math.Abs(req.DeliverE) > 0.5 {
		t.Errorf("Deliver N/E = %.2f/%.2f, want ~111.2/~0 (relative to pickup)", req.DeliverN, req.DeliverE)
	}
}

// TestDeliveryRequestJSONContract locks the JSON field names to the shared
// contract consumed by airpost_service.py.
func TestDeliveryRequestJSONContract(t *testing.T) {
	req := DeliveryRequest{
		OrderID: "AP-1", DroneID: 51, TakeoffID: 2, PickupID: 2,
		DeliverN: 75, DeliverE: 75, LandingID: 7, Cruise: 30,
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"order_id", "drone_id", "takeoff_id", "pickup_id", "deliver_N", "deliver_E", "landing_id", "cruise"} {
		if _, ok := got[key]; !ok {
			t.Errorf("missing contract field %q in %s", key, b)
		}
	}
}

func TestDeliveryStatusIsDelivered(t *testing.T) {
	tests := []struct {
		name   string
		status DeliveryStatus
		want   bool
	}{
		{"delivered state", DeliveryStatus{State: "delivered"}, true},
		{"terminal pass", DeliveryStatus{State: "done", Result: "PASS"}, true},
		{"enroute not delivered", DeliveryStatus{State: "enroute_delivery"}, false},
		{"failed", DeliveryStatus{State: "failed", Result: "FAIL"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsDelivered(); got != tt.want {
				t.Errorf("IsDelivered() = %v, want %v", got, tt.want)
			}
		})
	}
}
