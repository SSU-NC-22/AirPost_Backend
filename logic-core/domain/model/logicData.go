package model

import (
	"time"
)

type LogicData struct {
	Values     map[string]float64 `json:"values"` // sensor values
	NodeID     int                `json:"node_id"`
	Node       Node               `json:"node"`
	Timestamp  time.Time          `json:"timestamp"`
}

type Logic struct {
	ID			int       `json:"id"`
	LogicName	string    `json:"logic_name"`
	Elems		[]Element `json:"elems"`
	NodeID 		int       `json:"node_id"`
}

type Element struct {
	Elem string                 `json:"elem"`
	Arg  map[string]interface{} `json:"arg"`
}
