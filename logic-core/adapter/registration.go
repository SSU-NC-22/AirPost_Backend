package adapter

import (
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/eunnseo/AirPost/logic-core/domain/model"
)

// INIT
// adapter : []Sink + Sink.Nodes + Sink.Nodes.Sensors + Sink.Nodes.Sesnors.Logics + Sink.Nodes.Sensors.SensorValues
// action : create Nodes, Sensors, Logics

// DeleteSink
// adapter : []Node + Node.Sink (Node.Sensors X)
// action : delete Nodes

// CreateNode
// adapter : Node + Node.Sink + Node.Sensors.SensorValues + Node.Sensors.Logics
// action : create Node, Sensors, Logics

// adapter.Node -> model.Node + []adapter.Sensor
// adapter.Sensor -> model.Sensor + []adpater.Logic
// []adapter.Logic -> []model.Logic

// DeleteNode
// adapter : Node (Node.Sink, Node.Sensors X)
// action : delete Node

// DeleteSensor
// adapter : Sensor + Sensor.Logics
// action : delete Sensor, Logics

// CreateLogic
// adapter : Logic (Logic.Sensor X)
// action : create Logic

// DeleteLogic
// adapter : Logic (Logic.Sensor X)
// action : delete Logic

type Logic struct {
	ID		int    `json:"id"`
	Name	string `json:"name"`
	Elems	string `json:"elems"`
	NodeID	int    `json:"node_id"`
}

func LogicToModel(l *Logic) (model.Logic, error) {
	log.Println("LogicToModel")
	var elems []model.Element
	if err := json.Unmarshal([]byte(l.Elems), &elems); err != nil {
		return model.Logic{}, err
	} else {
		return model.Logic{
			ID:        l.ID,
			LogicName: l.Name,
			Elems:     elems,
			NodeID:    l.NodeID,
		}, nil
	}
}

func LogicsToModels(ll []Logic) []model.Logic {
	res := make([]model.Logic, 0, len(ll))
	for _, l := range ll {
		if ml, err := LogicToModel(&l); err == nil {
			res = append(res, ml)
		}
	}
	return res
}

type SensorValue struct {
	NodeID		int    `json:"node_id"`
	ValueName	string `json:"value_name"`
	Index		int    `json:"index"`
}

type Node struct {
	ID				int				`json:"id"`
	Name			string			`json:"name"`
	Type			string			`json:"type"`
	LocLat			float64			`json:"lat"`
	LocLon			float64			`json:"lng"`
	LocAlt			float64			`json:"alt"`
	SinkID			int				`json:"sink_id"`
	Sink			Sink			`json:"sink"`
	SensorValues	[]SensorValue	`json:"sensor_values"`
	Logics			[]Logic			`json:"logics"`
}

func NodeToModel(n *Node, sn string) (model.Node, []Logic) {
	sv := make([]string, len(n.SensorValues))
	for i, v := range n.SensorValues {
		sv[i] = v.ValueName
	}
	return model.Node{
		Name: n.Name,
		Type: n.Type,
		Location: model.Location{
			Lat: n.LocLat,
			Lon: n.LocLon,
			Alt: n.LocAlt,
		},
		SinkName:     sn,
		Sid:          n.SinkID,
		Nid:          n.ID,
		SensorValues: sv,
	}, n.Logics
}

type Sink struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Addr  string `json:"addr"`
	Nodes []Node `json:"nodes"`
}

type Topic struct {
	Name string `json:"name"`
}

type LogicService struct {
	Addr  string `json:"addr"`
	Topic Topic  `json:"topic"`
}

type Delivery struct {
	ID			  int		`json:"id"`
	OrderNum	  string	`json:"order_num"`
	DroneID		  int		`json:"drone_id"`
	Drone		  Node		`json:"drone_node"`

	SrcName		  string	`json:"src_name"`
	SrcPhone	  string	`json:"src_phone"`
	SrcStationID  int		`json:"src_station_id"`

	DestName	  string	`json:"dest_name"`
	DestPhone	  string	`json:"dest_phone"`
	DestStationID int		`json:"dest_station_id"`
	
	CreatedAt	  time.Time `json:"created_at"`
}

func DeliveryToModel(d *Delivery) (model.Delivery) {
	return model.Delivery{
		Did:		   d.ID,
		OrderNum:	   d.OrderNum,
		SrcName:	   d.SrcName,
		SrcPhone:	   d.SrcPhone,
		SrcStationID:  d.SrcStationID,
		DestName:	   d.DestName,
		DestPhone:     d.DestPhone,
		DestTagID: 	   d.DestStationID,
	}
}

type Path struct {
	StationID int	  `json:"station_id"`
	TagID 	  int	  `json:"tag_id"`
	Path   	  string  `json:"path"`
	Distance  float64 `json:"distance"`
}

func PathToModel(station *model.Node, tag *model.Node) (model.Path) {
	powLon := math.Pow((station.Location.Lon - tag.Location.Lon), 2)
	powLat := math.Pow((station.Location.Lat - tag.Location.Lat), 2)
	dist := math.Pow((powLon + powLat), 0.5)
	return model.Path{
		StationID: station.Nid,
		TagID:	   tag.Nid,
		Path:	   "",
		Distance:  dist,
	}
}