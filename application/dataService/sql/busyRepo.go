package sql

import (
	"time"

	"github.com/eunnseo/AirPost/application/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// busyRepo persists drone in-flight reservations in MySQL with a TTL (model.DroneBusy.BusyUntil),
// replacing the dispatcher's in-memory busy map. DB-backed so reservations survive an application
// restart; TTL so a sortie that never reports completion auto-frees its drone instead of leaking it.
type busyRepo struct{ db *gorm.DB }

func NewBusyRepo() *busyRepo {
	return &busyRepo{db: dbConn}
}

// SetBusy reserves a drone for an order until `until` (upsert on drone_id).
func (r *busyRepo) SetBusy(droneID int, orderNum string, until time.Time) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "drone_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"order_num", "busy_until"}),
	}).Create(&model.DroneBusy{DroneID: droneID, OrderNum: orderNum, BusyUntil: until}).Error
}

// FreeByOrder releases whatever drone the given order reserved.
func (r *busyRepo) FreeByOrder(orderNum string) error {
	return r.db.Where("order_num = ?", orderNum).Delete(&model.DroneBusy{}).Error
}

// IsBusy reports whether the drone has an UNEXPIRED reservation. Expired rows are ignored (the TTL),
// so a crashed/never-completed sortie stops blocking the drone once BusyUntil has passed.
func (r *busyRepo) IsBusy(droneID int) (bool, error) {
	var n int64
	err := r.db.Model(&model.DroneBusy{}).
		Where("drone_id = ? AND busy_until > ?", droneID, time.Now()).
		Count(&n).Error
	return n > 0, err
}
