package usecase

import (
	"github.com/eunnseo/AirPost/logic-core/adapter"
)

type EventUsecase interface {
	CreateSink(s *adapter.Sink) error
	DeleteSink(nl []adapter.Node) error
	
	CreateNode(n *adapter.Node, sn string) error
	DeleteNode(n *adapter.Node) error

	CreateLogic(l *adapter.Logic) error
	DeleteLogic(l *adapter.Logic) error

	CreateDelivery(d *adapter.Delivery) error
}

type LogicCoreUsecase interface {
	AppendSinkAddr(sa *adapter.SinkAddr) error
}
