package sql

import (
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

// Seed IDs. Station node IDs intentionally equal sim station IDs in
// simulation/tests/airpost_sites.json (the sim looks up world coords by these),
// so a fresh `compose up` flies a real sortie with no manual setup.
const (
	seedDroneID   = 50 // drone node
	seedTakeoffID = 1  // source station  -> airpost_sites station 1
	seedLandingID = 7  // landing station -> airpost_sites station 7
	seedTagID     = 30 // destination tag (drop point)

	sinkDrone   = 1 // drone-sink   (seeded in main.go)
	sinkStation = 2 // station-sink
	sinkTag     = 3 // tag-sink
)

// Seed inserts a usable demo topology (source station + landing station + drop
// tag + drone, plus the station-drone link and a path) on first run so the GUI
// order -> drone-flies -> track flow works end to end without any manual curl.
// It is idempotent (FirstOrCreate by primary key) and a no-op when SEED_DEMO=0.
func Seed() {
	if os.Getenv("SEED_DEMO") == "0" {
		return
	}

	// Base lat/lon for the source station (gives the UI map a real location).
	const baseLat, baseLon = 37.5000, 127.0000

	// Drop point 15 m north / 15 m east of the source station. The dispatcher
	// re-derives these offsets from the tag's lat/lon, so place the tag there.
	const deliverN, deliverE = 15.0, 15.0
	rad := math.Pi / 180
	tagLat := baseLat + (deliverN/earthRadiusMeters)/rad
	tagLon := baseLon + (deliverE/(earthRadiusMeters*math.Cos(baseLat*rad)))/rad

	// Landing station offset (purely for a distinct map marker / track endpoint).
	landLat := baseLat + (120.0/earthRadiusMeters)/rad
	landLon := baseLon + (40.0/(earthRadiusMeters*math.Cos(baseLat*rad)))/rad

	nodes := []model.Node{
		{ID: seedDroneID, Name: "drone-1", Type: "drone", LocLat: baseLat, LocLon: baseLon, LocAlt: 0, SinkID: sinkDrone},
		{ID: seedTakeoffID, Name: "station-src", Type: "station", LocLat: baseLat, LocLon: baseLon, LocAlt: 0, SinkID: sinkStation},
		{ID: seedLandingID, Name: "station-dest", Type: "station", LocLat: landLat, LocLon: landLon, LocAlt: 0, SinkID: sinkStation},
		{ID: seedTagID, Name: "tag-home", Type: "tag", LocLat: tagLat, LocLon: tagLon, LocAlt: 0, SinkID: sinkTag},
	}
	for i := range nodes {
		if err := dbConn.Clauses(clause.OnConflict{DoNothing: true}).
			Omit(clause.Associations).Create(&nodes[i]).Error; err != nil {
			log.Printf("seed: node %d: %v", nodes[i].ID, err)
		}
	}

	// The source station has a usable drone parked on it.
	sd := model.StationDrone{StationID: seedTakeoffID, DroneID: seedDroneID, Usable: true}
	if err := dbConn.Clauses(clause.OnConflict{DoNothing: true}).
		Omit(clause.Associations).Create(&sd).Error; err != nil {
		log.Printf("seed: station_drone: %v", err)
	}

	// One path tag -> source station so GetShortestPathStation resolves a landing
	// station for the drop tag. Landing back at the source station keeps the drone
	// parked there (Usable), so the demo is repeatable on the same `compose up`
	// (RegistDelivery only migrates the drone when src != dest station).
	path := model.Path{ID: 1, StationID: seedTakeoffID, TagID: seedTagID, Path: "[]", Distance: 10.0}
	if err := dbConn.Clauses(clause.OnConflict{DoNothing: true}).Create(&path).Error; err != nil {
		log.Printf("seed: path: %v", err)
	}

	log.Printf("seed: demo topology ready (src station %d, dest station %d, tag %d, drone %d)",
		seedTakeoffID, seedLandingID, seedTagID, seedDroneID)
}
