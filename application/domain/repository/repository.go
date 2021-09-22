package repository

import (
	"github.com/eunnseo/AirPost/application/adapter"
	"github.com/eunnseo/AirPost/application/domain/model"
)

type SinkRepo interface {
	GetPages(size int) int
	FindsWithTopic() ([]model.Sink, error)
	FindsPage(p adapter.Page) ([]model.Sink, error)
	FindsByTopicIDWithNodesSensorsValuesLogics(tid int) (sl []model.Sink, err error)
	FindByIDWithNodesSensorsValuesTopic(id int) (*model.Sink, error)
	Create(*model.Sink) error
	Delete(*model.Sink) error
}

type NodeRepo interface {
	GetPages(p adapter.Page) int
	FindsWithSensorsValues() ([]model.Node, error)
	FindsPage(p adapter.Page) (nl []model.Node, err error)
	FindsSquare(sq adapter.Square) (nl []model.Node, err error)
	FindsBySinkIDWithSensorValues(sinkid int) (nl []model.Node, err error)
	FindsByID(id int) (*model.Node, error)
	Create(*model.Node) error
	Delete(*model.Node) error
}

type LogicRepo interface {
	FindsWithNodeValues() ([]model.Logic, error)
	Create(*model.Logic) error
	Delete(*model.Logic) error
}

type LogicServiceRepo interface {
	Finds() ([]model.LogicService, error)
	FindsWithTopic() ([]model.LogicService, error)
	FindsByTopicID(TopicID int) ([]model.LogicService, error)
	FindByAddr(addr string) (l *model.LogicService, err error)
	Create(*model.LogicService) error
	Delete(*model.LogicService) error
}

type TopicRepo interface {
	FindsWithLogicService() ([]model.Topic, error)
	Create(*model.Topic) error
	Delete(*model.Topic) error
}

type DeliveryRepo interface {
	FindsByOrderNum(ordernum int) (dl model.Delivery, err error)
	Create(*model.Delivery) error
	Delete(*model.Delivery) error
}

type PathRepo interface {
	Finds() (pl []model.Path, err error)
	Create(*model.Path) error
	Delete(*model.Path) error
}

type StationDroneRepo interface {
	Find(stationid int, droneid int) (sd *model.StationDrone, err error)
	FindsByStationID(stationid int) ([]model.StationDrone, error)
	FindsByDroneID(droneid int) ([]model.StationDrone, error)
	Create(*model.StationDrone) error
	Delete(*model.StationDrone) error
	DeleteByStationID(*model.StationDrone) error
	DeleteByDroneID(*model.StationDrone) error
}