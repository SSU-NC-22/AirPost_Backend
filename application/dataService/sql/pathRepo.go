package sql

import (
	"github.com/eunnseo/AirPost/application/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pathRepo struct {
	db *gorm.DB
}

func NewPathRepo() *pathRepo {
	return &pathRepo{
		db: dbConn,
	}
}

func (ptr *pathRepo) Finds() (pl []model.Path, err error) {
	return pl, ptr.db.Find(&pl).Error
}

func (ptr *pathRepo) Create(p *model.Path) error {
	return ptr.db.Omit(clause.Associations).Create(p).Error
}

func (ptr *pathRepo) Delete(p *model.Path) error {
	return ptr.db.Delete(p).Error
}
