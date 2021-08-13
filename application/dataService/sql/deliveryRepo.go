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

func (dlr *deliveryRepo) Create(d *model.Delivery) error {
	return dlr.db.Omit(clause.Associations).Create(d).Error
}

func (dlr *deliveryRepo) Delete(d *model.Delivery) error {
	return dlr.db.Delete(d).Error
}
