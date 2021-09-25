package adapter

type States struct {
	Timestamp string     `json:"timestamp"`
	State     HealthInfo `json:"healthinfo"`
}

type NodeState struct {
	NodeID   int  	   `json:"nid"`
	State    bool	   `json:"state"`
	Battery  int 	   `json:"battery"`
	Location []float64 `json:"location"` // [lat, lon, alt]
}
type HealthInfo struct {
	SinkID int         `json:"sid"`
	State  []NodeState `json:"state"`
}
