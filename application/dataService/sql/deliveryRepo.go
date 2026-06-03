package sql

import (
	"github.com/eunnseo/AirPost/application/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type deliveryRepo struct {
	db *gorm.DB
}

func NewDeliveryRepo() *deliveryRepo {
	return &deliveryRepo{
		db: dbConn,
	}
}

func (dlr *deliveryRepo) FindsByOrderNum(ordernum string) (dl model.Delivery, err error) {
	return dl, dlr.db.Where("order_num=?", ordernum).Find(&dl).Error
}

func (dlr *deliveryRepo) Create(d *model.Delivery) error {
	return dlr.db.Omit(clause.Associations).Create(d).Error
}

func (dlr *deliveryRepo) Delete(d *model.Delivery) error {
	return dlr.db.Delete(d).Error
}

// DeleteByDroneID removes every delivery booked to a drone. Deliveries reference the drone node
// without an ON DELETE cascade, so they must be cleared before the drone node itself can be deleted
// (e.g. when an offline drone is pruned from the IoT fleet).
func (dlr *deliveryRepo) DeleteByDroneID(droneid int) error {
	return dlr.db.Where("drone_id=?", droneid).Delete(&model.Delivery{}).Error
}
