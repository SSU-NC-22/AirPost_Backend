package eventUC

import (
	"log"

	"github.com/eunnseo/AirPost/logic-core/adapter"
	"github.com/eunnseo/AirPost/logic-core/domain/repository"
	"github.com/eunnseo/AirPost/logic-core/domain/service"
)

type eventUsecase struct {
	rr repository.RegistRepo
	ls service.LogicService
}

func NewEventUsecase(rr repository.RegistRepo, ls service.LogicService) *eventUsecase {
	return &eventUsecase{
		rr: rr,
		ls: ls,
	}
}

func (eu *eventUsecase) DeleteSink(nl []adapter.Node) error {
	for _, n := range nl {
		eu.rr.DeleteNode(n.ID)
	}
	return nil
}

func (eu *eventUsecase) CreateNode(n *adapter.Node, sn string) error {
	// node
	mn, all := adapter.NodeToModel(n, sn)
	eu.rr.CreateNode(n.ID, &mn) // regist in nodeRepo

	// all := []adapter.Logic{}
	// for _, as := range asl {
	// 	ms, tempAll := adapter.SensorToModel(&as)
	// 	all = append(all, tempAll...)
	// 	eu.rr.CreateSensor(as.ID, &ms)
	// }

	// logic
	mll := adapter.LogicsToModels(all)
	for _, ml := range mll {
		eu.ls.CreateAndStartLogic(&ml)
	}

	return nil
}

func (eu *eventUsecase) DeleteNode(n *adapter.Node) error {
	return eu.rr.DeleteNode(n.ID)
}

func (eu *eventUsecase) CreateLogic(l *adapter.Logic) error {
	log.Println("in eu.CreateLogic")
	if ml, err := adapter.LogicToModel(l); err != nil {
		return err
	} else {
		log.Println("in eu.CreateLogic.good")
		return eu.ls.CreateAndStartLogic(&ml)
	}
}

func (eu *eventUsecase) DeleteLogic(l *adapter.Logic) error {
	return eu.ls.RemoveLogic(l.NodeID, l.ID)
}

func (eu *eventUsecase) CreateDelivery(d *adapter.Delivery) error {
	log.Println("in eu.CreateDelivery")
	
	return nil
}
