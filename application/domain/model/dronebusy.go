package model

import "time"

// DroneBusy is a drone's in-flight reservation, persisted so it survives an application restart and
// is shared correctly even if dispatch ran on a previous process. BusyUntil is a TTL: a reservation
// is only "busy" while BusyUntil is in the future, so a sortie whose completion status never arrives
// (e.g. the drone/sim crashed) auto-frees the drone after the TTL instead of leaking it forever
// (which used to wedge dispatch into "no usable drone available").
type DroneBusy struct {
	DroneID   int       `json:"drone_id" gorm:"primaryKey;autoIncrement:false"`
	OrderNum  string    `json:"order_num" gorm:"type:varchar(32);index"`
	BusyUntil time.Time `json:"busy_until" gorm:"not null"`
}

func (DroneBusy) TableName() string {
	return "drone_busy"
}
