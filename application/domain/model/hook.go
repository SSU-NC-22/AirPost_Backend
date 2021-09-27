package model

import (
	"fmt"

	"gorm.io/gorm"
)

var orderByASC = func(db *gorm.DB) *gorm.DB {
	return db.Order("sensor_values.index ASC")
}

// sink
func (s *Sink) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Preload("Topic").Find(s).Error
}

func (s *Sink) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Preload("Topic").Preload("Nodes").Find(s).Error
}

// node
func (n *Node) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Preload("Sink.Topic").Preload("Sink").Preload("Logics").Preload("SensorValues", orderByASC).Find(n).Error
}

func (n *Node) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Preload("Sink.Topic").Preload("Sink").Find(n).Error
}

// logic
func (l *Logic) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Preload("Node").Find(l).Error
}

func (l *Logic) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Preload("Node").Find(l).Error
}

// logicService
func (l *LogicService) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Preload("Topic.Sinks.Nodes.Logics").Preload("Topic.Sinks.Nodes.SensorValues", orderByASC).Preload("Topic.Sinks.Nodes").Preload("Topic.Sinks").Preload("Topic").Find(l).Error
}

func (l *LogicService) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Preload("Topic").Find(l).Error
}

// topic
func (t *Topic) BeforeDelete(tx *gorm.DB) (err error) {
	if err = tx.Preload("LogicServices").Find(t).Error; err != nil {
		return err
	}
	if len(t.LogicServices) != 0 {
		return fmt.Errorf("there are logic-services that consume topic : %s", t.Name)
	}
	return nil
}
