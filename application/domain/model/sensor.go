package model

// modified
type SensorValue struct {
	// SensorID  int    `json:"sensor_id" gorm:"primaryKey"`
	NodeID		int		`json:"node_id" gorm:"primaryKey"`
	ValueName	string	`json:"value_name" gorm:"primaryKey;type:varchar(32)"`
	Index		int		`json:"index" gorm:"not null"`
}

func (SensorValue) TableName() string {
	return "sensor_values"
}
