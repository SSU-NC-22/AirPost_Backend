package model

type Node struct {
	Name			string		`json:"name"`
	Type			string		`json:"type"`
	Location		Location	`json:"location"`
	SinkName		string		`json:"sink_name"`
	Sid				int			`json:"sid"`	// sink id
	Nid				int			`json:"nid"`	// node id
	SensorValues	[]string	`json:"sensor_values"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Alt float64 `json:"alt"`
}

type Sink struct {
	Addr string `json:"addr"`
}
type Nodeinfo struct {
	SinkID int `json:"sink_id"`
}

type Delivery struct {
	Did			  int  	 `json:"did"`
	OrderNum	  string `json:"order_num"`
	Email         string `json:"email"`
	SrcName		  string `json:"src_name"`
	SrcPhone	  string `json:"src_phone"`
	SrcStationID  int	 `json:"src_station_id"`
	DestName	  string `json:"dest_name"`
	DestPhone	  string `json:"dest_phone"`
	DestTagID     int	 `json:"dest_tag_id"`
}

type Path struct {
	Pid		  int	  `json:"pid"`
	StationID int	  `json:"station_id"`
	TagID 	  int	  `json:"tag_id"`
	Path   	  string  `json:"path"`
	Distance  float64 `json:"distance"`
}
