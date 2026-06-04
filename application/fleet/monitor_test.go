package fleet

import (
	"testing"
	"time"
)

// nearby lat/lon offsets around the origin that resolve to small metre distances.
func atMetres(north, east float64) (lat, lon float64) {
	lat = originLat + north/earthR*180.0/3.141592653589793
	lon = originLon + east/(earthR*0.7933533402912352)*180.0/3.141592653589793 // cos(37.5deg)
	return
}

func has(advs []Advisory, id int, hold bool) bool {
	for _, a := range advs {
		if a.DroneID == id && a.Hold == hold {
			return true
		}
	}
	return false
}

func TestHoldsHigherIdOnConflict(t *testing.T) {
	m := NewMonitor()
	now := time.Unix(1000, 0)
	la1, lo1 := atMetres(0, 0)
	la2, lo2 := atMetres(3, 0) // 3 m north of drone 51 — well within minSepXY, same altitude
	m.Update(51, la1, lo1, 30, now)
	m.Update(52, la2, lo2, 30, now)

	advs := m.Deconflict(now)
	if !has(advs, 52, true) {
		t.Fatalf("expected drone 52 (higher id) to be held, got %+v", advs)
	}
	if has(advs, 51, true) {
		t.Fatalf("drone 51 (lower id, priority) must not be held, got %+v", advs)
	}
}

func TestClearsWhenSeparated(t *testing.T) {
	m := NewMonitor()
	now := time.Unix(1000, 0)
	la1, lo1 := atMetres(0, 0)
	la2, lo2 := atMetres(3, 0)
	m.Update(51, la1, lo1, 30, now)
	m.Update(52, la2, lo2, 30, now)
	_ = m.Deconflict(now) // 52 held

	// 52 moves 50 m away -> conflict resolved -> 52 cleared
	la2, lo2 = atMetres(50, 0)
	m.Update(52, la2, lo2, 30, now)
	advs := m.Deconflict(now)
	if !has(advs, 52, false) {
		t.Fatalf("expected drone 52 to be cleared once separated, got %+v", advs)
	}
}

func TestVerticalSeparationAvoidsConflict(t *testing.T) {
	m := NewMonitor()
	now := time.Unix(1000, 0)
	la1, lo1 := atMetres(0, 0)
	la2, lo2 := atMetres(3, 0) // horizontally close...
	m.Update(51, la1, lo1, 30, now)
	m.Update(52, la2, lo2, 40, now) // ...but 10 m higher (different band) -> no conflict
	if advs := m.Deconflict(now); len(advs) != 0 {
		t.Fatalf("vertically separated drones must not conflict, got %+v", advs)
	}
}

func TestStalePositionIgnored(t *testing.T) {
	m := NewMonitor()
	t0 := time.Unix(1000, 0)
	la1, lo1 := atMetres(0, 0)
	la2, lo2 := atMetres(3, 0)
	m.Update(51, la1, lo1, 30, t0)
	m.Update(52, la2, lo2, 30, t0)
	// evaluate 10 s later: both positions are stale -> ignored -> no advisories
	if advs := m.Deconflict(t0.Add(10 * time.Second)); len(advs) != 0 {
		t.Fatalf("stale positions must be ignored, got %+v", advs)
	}
}

func TestNoFlapWhileConflictPersists(t *testing.T) {
	m := NewMonitor()
	now := time.Unix(1000, 0)
	la1, lo1 := atMetres(0, 0)
	la2, lo2 := atMetres(3, 0)
	m.Update(51, la1, lo1, 30, now)
	m.Update(52, la2, lo2, 30, now)
	first := m.Deconflict(now)
	if !has(first, 52, true) {
		t.Fatalf("expected initial hold for 52, got %+v", first)
	}
	// same picture again -> no new advisory (hold state unchanged)
	if again := m.Deconflict(now); len(again) != 0 {
		t.Fatalf("hold must not re-fire while unchanged, got %+v", again)
	}
}
