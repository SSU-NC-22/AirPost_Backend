// Package fleet is the backend "control tower": it ingests every drone's live telemetry, keeps a
// picture of where the whole fleet is, and issues deconfliction advisories so two drones never get
// dangerously close. It complements the dispatcher's altitude-band assignment (vertical separation at
// dispatch time) with continuous, position-based monitoring — and tells a drone to HOLD when a
// conflict develops, CLEARing it once the airspace is free again.
//
// Division of responsibility (matches the on-drone ROS 2 node):
//   - backend (here): fleet-level planning + monitoring + hold/clear advisories.
//   - each drone (airpost_drone drone_node): executes its route and does its OWN local obstacle
//     avoidance. The backend never flies the drone; it only deconflicts the fleet.
package fleet

import (
	"math"
	"sort"
	"sync"
	"time"
)

// Geo origin shared with seed.go / the simulation, so a drone's reported lat/lon maps to the same
// local east/north metric frame the rest of the system uses.
const (
	originLat = 37.5
	originLon = 127.0
	earthR    = 6371000.0
)

// Separation thresholds: a conflict is two drones within minSepXY horizontally AND minSepZ
// vertically. The dispatcher's altitude bands (>= bandGap apart) normally keep drones out of each
// other's vertical slice; this catches the moments they converge — e.g. both descending toward the
// same station — and holds one off.
const (
	defaultSepXY = 12.0            // metres
	defaultSepZ  = 4.0             // metres
	defaultStale = 5 * time.Second // ignore positions older than this (drone offline / no telemetry)
)

// Position is a drone's last known location in the local NED-ish metric frame (north, east, up).
type Position struct {
	N, E, Alt float64
	Updated   time.Time
}

// Advisory tells one drone to hold (Hold=true) or that it is clear to proceed (Hold=false). The
// monitor emits an advisory only when a drone's hold state CHANGES, so the transport stays quiet.
type Advisory struct {
	DroneID int
	Hold    bool
	Reason  string
}

// Monitor is the live fleet picture + deconfliction state. Safe for concurrent use.
type Monitor struct {
	mu     sync.Mutex
	pos    map[int]Position
	held   map[int]bool
	sepXY  float64
	sepZ   float64
	stale  time.Duration
}

// NewMonitor builds a Monitor with the default separation/staleness thresholds.
func NewMonitor() *Monitor {
	return &Monitor{
		pos:   make(map[int]Position),
		held:  make(map[int]bool),
		sepXY: defaultSepXY,
		sepZ:  defaultSepZ,
		stale: defaultStale,
	}
}

// Update records a drone's latest position from its telemetry (lat/lon degrees, altitude metres).
func (m *Monitor) Update(droneID int, lat, lon, alt float64, now time.Time) {
	n, e := localNE(lat, lon)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pos[droneID] = Position{N: n, E: e, Alt: alt, Updated: now}
}

// Deconflict scans the fresh fleet picture and returns advisories for any drone whose hold state
// changed. For each conflicting pair the HIGHER drone id is held (deterministic priority, so the same
// drone yields every tick and the decision never flaps); the other keeps flying. A drone with no
// remaining conflict that was previously held is cleared.
func (m *Monitor) Deconflict(now time.Time) []Advisory {
	m.mu.Lock()
	defer m.mu.Unlock()

	// fresh, sorted ids for deterministic pairing
	ids := make([]int, 0, len(m.pos))
	for id, p := range m.pos {
		if now.Sub(p.Updated) <= m.stale {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)

	conflict := make(map[int]bool, len(ids))
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			a, b := m.pos[ids[i]], m.pos[ids[j]]
			if math.Hypot(a.N-b.N, a.E-b.E) < m.sepXY && math.Abs(a.Alt-b.Alt) < m.sepZ {
				conflict[ids[j]] = true // hold the higher id; the lower keeps priority
			}
		}
	}

	var out []Advisory
	for _, id := range ids {
		want := conflict[id]
		if want != m.held[id] {
			m.held[id] = want
			reason := ""
			if want {
				reason = "fleet separation"
			}
			out = append(out, Advisory{DroneID: id, Hold: want, Reason: reason})
		}
	}
	return out
}

// Run starts the deconfliction loop: every `interval` it re-evaluates the fleet picture and calls
// `publish` for each drone whose hold state changed. Returns a stop function. `now` is injectable for
// tests; pass time.Now in production.
func (m *Monitor) Run(interval time.Duration, now func() time.Time, publish func(Advisory)) (stop func()) {
	ticker := time.NewTicker(interval)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				for _, a := range m.Deconflict(now()) {
					publish(a)
				}
			}
		}
	}()
	return func() { close(done) }
}

// localNE converts lat/lon (degrees) to local north/east metres about the shared geo origin
// (equirectangular — accurate at the few-kilometre scale of a delivery area).
func localNE(lat, lon float64) (north, east float64) {
	dLat := (lat - originLat) * math.Pi / 180.0
	dLon := (lon - originLon) * math.Pi / 180.0
	north = dLat * earthR
	east = dLon * earthR * math.Cos(originLat*math.Pi/180.0)
	return north, east
}
