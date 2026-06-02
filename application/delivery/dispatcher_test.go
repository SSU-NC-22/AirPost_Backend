package delivery

import (
	"testing"

	deliverymqtt "github.com/eunnseo/AirPost/application/delivery/mqtt"
)

// stubPublisher records published requests for assertions.
type stubPublisher struct{ reqs []deliverymqtt.DeliveryRequest }

func (s *stubPublisher) PublishRequest(r deliverymqtt.DeliveryRequest) error {
	s.reqs = append(s.reqs, r)
	return nil
}

func newTestDispatcher() *Dispatcher {
	return &Dispatcher{
		publisher:  &stubPublisher{},
		bands:      make(map[string]int),
		orderDrone: make(map[string]int),
		busy:       make(map[int]bool),
		notified:   make(map[string]bool),
	}
}

// TestBandsAreDistinctAndReused checks the control-tower deconfliction: concurrent
// missions get distinct cruise altitudes, and a band is reused once its mission
// lands (terminal status frees it).
func TestBandsAreDistinctAndReused(t *testing.T) {
	d := newTestDispatcher()

	a := d.reserveBand("A")
	b := d.reserveBand("B")
	c := d.reserveBand("C")
	if a == b || b == c || a == c {
		t.Fatalf("bands not distinct: A=%v B=%v C=%v", a, b, c)
	}
	if a != defaultCruiseAltitude {
		t.Errorf("first band = %v, want base %v", a, defaultCruiseAltitude)
	}
	if b != defaultCruiseAltitude+bandGap || c != defaultCruiseAltitude+2*bandGap {
		t.Errorf("bands not stacked by gap: B=%v C=%v", b, c)
	}

	// Land B (terminal status frees its band); it must be handed to the next
	// mission (lowest free). FAIL is terminal too and skips the email path.
	d.HandleStatus(deliverymqtt.DeliveryStatus{OrderID: "B", State: "aborted", Result: "FAIL"})
	if got := d.reserveBand("D"); got != b {
		t.Errorf("freed band not reused: D=%v, want %v", got, b)
	}
}

// TestReserveBandIdempotent ensures re-reserving the same order keeps its band.
func TestReserveBandIdempotent(t *testing.T) {
	d := newTestDispatcher()
	first := d.reserveBand("X")
	if again := d.reserveBand("X"); again != first {
		t.Errorf("re-reserve changed band: %v != %v", again, first)
	}
}
