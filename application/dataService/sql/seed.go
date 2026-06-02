package sql

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/eunnseo/AirPost/application/domain/model"
	"gorm.io/gorm/clause"
)

// earthRadiusMeters mirrors the projection used by the delivery mapping, so a
// tag placed a few meters from a station produces the same north/east offset the
// dispatcher publishes to the sim.
const earthRadiusMeters = 6371000.0

// Geo origin of the local Gazebo world (E=0,N=0). Station/tag lat/lon are derived
// from the sim's local E/N via this origin, so the backend's geometry (nearest
// station, ferry distances) matches what the simulator actually flies.
const (
	originLat = 37.5
	originLon = 127.0
)

// Sink IDs (seeded in main.go before Seed runs).
const (
	sinkDrone   = 1
	sinkStation = 2
	sinkTag     = 3
)

// ID ranges for the seeded demo fleet. Station IDs intentionally equal the sim
// station IDs in simulation/tests/airpost_sites.json (the simulator looks up each
// drone's spawn world-coords by station ID), so a fresh `compose up` flies real
// multi-drone sorties with no manual setup.
const (
	fleetSize     = 8  // stations == drones == tags
	droneIDBase   = 50 // drone i -> node id 50+i  (i = 1..fleetSize)
	tagIDBase     = 30 // tag   i -> node id 30+i
	tagDropNorthM = 12 // each tag sits this many meters north of its home station
)

// Sensor schemas per node type: the ordered value names a node's telemetry array
// carries (index = position). logic-core maps an incoming values[] onto these names
// before indexing into Elasticsearch, so telemetry must be published in this order.
// They mirror the field order documented in logic-core/kafkaProducer/producer.py.
var sensorSchema = map[string][]string{
	"drone":   {"lat", "long", "alt", "velocity", "batteryPer", "done"},
	"station": {"temperature", "humidity", "light", "lat", "long", "alt"},
	"tag":     {"lat", "long", "alt"},
}

// stationEN are the local east/north meters of stations 1..fleetSize, copied from
// simulation/tests/airpost_sites.json (validated clear helipads, well separated so
// the drones spawn at distinct positions). Index 0 is station ID 1.
var stationEN = [fleetSize][2]float64{
	{320, -32}, {352, 64}, {384, 32}, {32, -80},
	{176, 256}, {224, -32}, {-144, 176}, {-144, 80},
}

// enToLatLon projects local east/north meters to lat/lon using the world origin
// (the inverse of the equirectangular projection in delivery/mqtt/mapping.go).
func enToLatLon(east, north float64) (lat, lon float64) {
	rad := math.Pi / 180
	lat = originLat + (north/earthRadiusMeters)/rad
	lon = originLon + (east/(earthRadiusMeters*math.Cos(originLat*rad)))/rad
	return lat, lon
}

// haversineMeters is the great-circle distance between two lat/lon points; used to
// fill Path.Distance so GetShortestPathStation resolves the genuinely nearest
// landing station for each drop tag.
func haversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	rad := math.Pi / 180
	dLat := (lat2 - lat1) * rad
	dLon := (lon2 - lon1) * rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return earthRadiusMeters * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// Seed inserts a usable multi-drone demo topology on first run: fleetSize stations
// (each a sim helipad with one usable drone parked on it), one drop tag per station,
// and a full tag->station path matrix so the dispatcher can pick the nearest landing
// station. It is idempotent (OnConflict DoNothing) and a no-op when SEED_DEMO=0.
func Seed() {
	if os.Getenv("SEED_DEMO") == "0" {
		return
	}

	var (
		nodes   []model.Node
		sensors []model.SensorValue
		links   []model.StationDrone
		paths   []model.Path
		pathID  = 1
		stLat   [fleetSize]float64
		stLon   [fleetSize]float64
		tagLat  [fleetSize]float64
		tagLon  [fleetSize]float64
	)

	// addSensors records a node's telemetry value schema (ordered by index).
	addSensors := func(nodeID int, kind string) {
		for i, name := range sensorSchema[kind] {
			sensors = append(sensors, model.SensorValue{NodeID: nodeID, ValueName: name, Index: i})
		}
	}

	for i := 0; i < fleetSize; i++ {
		id := i + 1
		lat, lon := enToLatLon(stationEN[i][0], stationEN[i][1])
		stLat[i], stLon[i] = lat, lon
		// Drop tag a few meters north of its home station.
		tLat, tLon := enToLatLon(stationEN[i][0], stationEN[i][1]+tagDropNorthM)
		tagLat[i], tagLon[i] = tLat, tLon

		droneID := droneIDBase + id
		tagID := tagIDBase + id

		nodes = append(nodes,
			model.Node{ID: id, Name: fmt.Sprintf("station-%d", id), Type: "station",
				LocLat: lat, LocLon: lon, SinkID: sinkStation},
			model.Node{ID: droneID, Name: fmt.Sprintf("drone-%d", id), Type: "drone",
				LocLat: lat, LocLon: lon, SinkID: sinkDrone},
			model.Node{ID: tagID, Name: fmt.Sprintf("tag-%d", id), Type: "tag",
				LocLat: tLat, LocLon: tLon, SinkID: sinkTag},
		)
		addSensors(id, "station")
		addSensors(droneID, "drone")
		addSensors(tagID, "tag")

		// One usable drone parked on each station.
		links = append(links, model.StationDrone{StationID: id, DroneID: droneID, Usable: true})
	}

	// Full tag -> station path matrix with real distances, so landing resolves to
	// the station nearest the drop point (which, for a cross-map delivery, is the
	// destination station rather than the takeoff one).
	for t := 0; t < fleetSize; t++ {
		tagID := tagIDBase + (t + 1)
		for s := 0; s < fleetSize; s++ {
			stationID := s + 1
			dist := haversineMeters(tagLat[t], tagLon[t], stLat[s], stLon[s])
			paths = append(paths, model.Path{
				ID: pathID, StationID: stationID, TagID: tagID, Path: "[]", Distance: dist,
			})
			pathID++
		}
	}

	create := func(what string, v interface{}) {
		if err := dbConn.Clauses(clause.OnConflict{DoNothing: true}).
			Omit(clause.Associations).Create(v).Error; err != nil {
			log.Printf("seed: %s: %v", what, err)
		}
	}
	create("nodes", &nodes)
	create("sensor_values", &sensors)
	create("station_drone", &links)
	create("paths", &paths)

	log.Printf("seed: demo topology ready (%d stations, %d drones, %d tags, %d paths, %d sensor values)",
		fleetSize, fleetSize, fleetSize, len(paths), len(sensors))
}
