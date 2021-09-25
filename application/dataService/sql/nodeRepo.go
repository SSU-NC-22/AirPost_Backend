package sql

import (
	"github.com/eunnseo/AirPost/application/adapter"
	"github.com/eunnseo/AirPost/application/domain/model"
	"gorm.io/gorm"
)

var orderByASC = func(db *gorm.DB) *gorm.DB {
	return db.Order("sensor_values.index ASC")
}

type nodeRepo struct {
	db *gorm.DB
}

func NewNodeRepo() *nodeRepo {
	return &nodeRepo{
		db: dbConn,
	}
}

func (ndr *nodeRepo) GetPages(p adapter.Page) int {
	temp := []model.Node{}
	if p.Sink != 0 {
		result := ndr.db.Where("sink_id=?", p.Sink).Find(&temp)
		count := int(result.RowsAffected)
		return (count / p.Size) + 1
	} else {
		result := ndr.db.Find(&temp)
		count := int(result.RowsAffected)
		return (count / p.Size) + 1
	}
}

func (ndr *nodeRepo) FindsWithSensorsValues() (nl []model.Node, err error) {
	return nl, ndr.db.Preload("SensorValues", orderByASC).Find(&nl).Error
}

func (ndr *nodeRepo) FindsPage(p adapter.Page) (nl []model.Node, err error) {
	offset := p.GetOffset()
	if p.Sink == 0 {
		return nl, ndr.db.Offset(offset).Limit(p.Size).Preload("SensorValues", orderByASC).Find(&nl).Error
	} else {
		return nl, ndr.db.Where("sink_id=?", p.Sink).Offset(offset).Limit(p.Size).Preload("SensorValues", orderByASC).Find(&nl).Error
	}
}

func (ndr *nodeRepo) FindsSquare(sq adapter.Square) (nl []model.Node, err error) {
	return nl, ndr.db.Where("loc_lon BETWEEN ? AND ?", sq.Left, sq.Right).Where("loc_lat BETWEEN ? AND ?", sq.Down, sq.Up).Preload("SensorValues", orderByASC).Find(&nl).Error
}

func (ndr *nodeRepo) FindsBySinkIDWithSensorValues(sinkid int) (nl []model.Node, err error) {
	return nl, ndr.db.Where("sink_id=?", sinkid).Preload("SensorValues", orderByASC).Find(&nl).Error
}

func (ndr *nodeRepo) FindsByID(id int) (*model.Node, error) {
	n := &model.Node{}
	return n, ndr.db.Where("id=?", id).Omit("SensorValues").Omit("Logics").Find(n).Error
}

func (ndr *nodeRepo) Create(n *model.Node) error {
	return ndr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("SensorValues").Omit("Logics").Create(n).Error; err != nil {
			return err
		}
		sv := n.SensorValues
		n.SensorValues = []model.SensorValue{}
		if err := tx.Model(n).Association("SensorValues").Append(sv); err != nil {
			return err
		}
		return nil
	})
}

func (ndr *nodeRepo) Delete(n *model.Node) error {
	return ndr.db.Delete(n).Error
}
