package model

import "time"

type Delivery struct {
	ID			  int		`json:"id" gorm:"primaryKey"`
	OrderNum	  string	`json:"order_num" gorm:"type:varchar(32);not null"`
	DroneID		  int		`json:"drone_id" gorm:"not null"`
	Drone		  Node		`json:"drone_node" gorm:"foreignKey:DroneID"`

	SrcName		  string	`json:"src_name" gorm:"type:varchar(32);not null"`
	SrcPhone	  string	`json:"src_phone" gorm:"type:varchar(32);not null"`
	SrcStationID  int		`json:"src_station_id" gorm:"not null"`

	DestName	  string	`json:"dest_name" gorm:"type:varchar(32);not null"`
	DestPhone	  string	`json:"dest_phone" gorm:"type:varchar(32);not null"`
	DestStationID int		`json:"dest_station_id" gorm:"not null"`
	
	CreatedAt	  time.Time `json:"created_at" gorm:"not null"`
}

func (Delivery) TableName() string {
	return "deliveries"
}
