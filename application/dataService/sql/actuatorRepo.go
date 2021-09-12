package sql

import (
	"github.com/eunnseo/AirPost/application/adapter"
	"github.com/eunnseo/AirPost/application/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type actuatorRepo struct {
	db *gorm.DB
}

func NewActuatorkRepo() *actuatorRepo {
	return &actuatorRepo{
		db: dbConn,
	}
}

func (acr *actuatorRepo) GetPages(size int) int {
	temp := []model.Sink{}
	result := acr.db.Find(&temp)
	count := int(result.RowsAffected)
	return (count / size) + 1
}

func (acr *actuatorRepo) FindsWithName() (al []model.Actuator, err error) {
	return al, acr.db.Find(&al).Error
}
func (acr *actuatorRepo) FindsPage(p adapter.Page) (al []model.Actuator, err error) {
	offset := p.GetOffset()
	return al, acr.db.Offset(offset).Limit(p.Size).Find(&al).Error
}

func (acr *actuatorRepo) Create(a *model.Actuator) error {
	return acr.db.Omit(clause.Associations).Create(a).Error
}

func (acr *actuatorRepo) Delete(a *model.Actuator) error {
	return acr.db.Delete(a).Error
}
