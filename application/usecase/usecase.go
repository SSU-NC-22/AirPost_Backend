package usecase

import (
	"github.com/eunnseo/AirPost/application/adapter"
	"github.com/eunnseo/AirPost/application/domain/model"
)

// for ui registration
type RegistUsecase interface {
	GetSinkPageCount(size int) int
	GetSinks() ([]model.Sink, error)
	GetSinksPage(p adapter.Page) ([]model.Sink, error)
	GetSinksByTopicID(tid int) ([]model.Sink, error)
	GetSinkByID(sid int) (*model.Sink, error)
	RegistSink(s *model.Sink) error
	UnregistSink(s *model.Sink) error

	GetNodePageCount(p adapter.Page) int
	GetNodes() ([]model.Node, error)
	GetNodesPage(p adapter.Page) ([]model.Node, error)
	GetNodesSquare(sq adapter.Square) ([]model.Node, error)
	GetNodesBySinkID(sinkid int) ([]model.Node, error)
	GetNodeByID(id int) (*model.Node, error)
	RegistNode(n *model.Node) error
	UnregistNode(n *model.Node) error

	GetActuatorPageCount(size int) int
	GetActuators() ([]model.Actuator, error)
	GetActuatorsPage(p adapter.Page) ([]model.Actuator, error)
	RegistActuator(a *model.Actuator) error
	UnregistActuator(a *model.Actuator) error

	GetLogics() ([]model.Logic, error)
	RegistLogic(l *model.Logic) error
	UnregistLogic(l *model.Logic) error

	GetLogicServices() ([]model.LogicService, error)
	UnregistLogicService(l *model.LogicService) error

	GetTopics() ([]model.Topic, error)
	RegistTopic(t *model.Topic) error
	UnregistTopic(t *model.Topic) error

	GetDeliveryByOrderNum(on int) (model.Delivery, error)
	RegistDelivery(d *model.Delivery) error
	UnregistDelivery(d *model.Delivery) error

	RegistPath(p *model.Path) error
	UnregistPath(p *model.Path) error

	GetStationDroneByStationID(stationid int) ([]model.StationDrone, error)
	RegistStationDrone(sd *model.StationDrone) error
	UnregistStationDrone(sd *model.StationDrone) error
}

// for event channel
type EventUsecase interface {
	RegistLogicService(l *model.LogicService) error
	CheckAndUnregistLogicServices() error

	CreateSinkEvent(s *model.Sink) error
	DeleteSinkEvent(s *model.Sink) error
	CreateNodeEvent(n *model.Node) error
	DeleteNodeEvent(n *model.Node) error
	CreateLogicEvent(l *model.Logic) error
	DeleteLogicEvent(l *model.Logic) error

	PostToSink(sid int) error

	CreateDeliveryEvent(d *model.Delivery) error
}
