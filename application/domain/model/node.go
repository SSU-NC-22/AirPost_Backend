package model

type Sink struct {
	ID      int    `json:"id" gorm:"primaryKey"`
	Name    string `json:"name" gorm:"type:varchar(32);unique;not null"`
	Addr    string `json:"addr" gorm:"type:varchar(32);not null"`
	TopicID int    `json:"topic_id" gorm:"not null"`
	Topic   Topic  `json:"topic" gorm:"foreignKey:TopicID"`
	Nodes   []Node `json:"nodes" gorm:"foreignKey:SinkID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}


func (Sink) TableName() string {
	return "sinks"
}

type Node struct {
	ID				int				`json:"id" gorm:"primaryKey"`
	Name			string			`json:"name" gorm:"type:varchar(32);unique;not null"`
	Type			string			`json:"type" gorm:"type:varchar(32)"`
	LocLat			float64			`json:"lat"`
	LocLon			float64			`json:"lng"`
	LocAlt			float64			`json:"alt"`
	SinkID			int				`json:"sink_id" gorm:"not null"`
	Sink			Sink			`json:"sink" gorm:"foreignKey:SinkID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SensorValues	[]SensorValue	`json:"sensor_values" gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Logics			[]Logic			`json:"logics" gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StationDrones	[]StationDrone	`json:"station_drone" gorm:"foreignKey:StationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Node) TableName() string {
	return "nodes"
}

type SensorValue struct {
	NodeID		int		`json:"node_id" gorm:"primaryKey"`
	ValueName	string	`json:"value_name" gorm:"primaryKey;type:varchar(32)"`
	Index		int		`json:"index" gorm:"not null"`
}

func (SensorValue) TableName() string {
	return "sensor_values"
}
