package eventUC

import (
	"log"

	"github.com/eunnseo/AirPost/logic-core/adapter"
	"github.com/eunnseo/AirPost/logic-core/domain/repository"
	"github.com/eunnseo/AirPost/logic-core/domain/service"
)

const (
	DRONE   = 1 // drone sink id
	STATION = 2 // station sink id
	TAG 	= 3 // tag sink id
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

func (eu *eventUsecase) CreateSink(s *adapter.Sink) error {
	err := eu.rr.AppendSinkAddr(s.ID, &s.Addr)
	if err != nil {
		log.Println("in eu.CreateSink, AppendSinkAddr error")
		return err
	}
	return nil
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
	eu.rr.CreateNode(n.ID, &mn)

	// logic
	mll := adapter.LogicsToModels(all)
	for _, ml := range mll {
		eu.ls.CreateAndStartLogic(&ml)
	}

	// path
	if mn.Type == "STA" {
		tags, _ := eu.rr.FindNodesBySinkID(TAG)
		for _, tag := range tags {
			log.Println("tag : ", tag.Name)

			path := adapter.PathToModel(&mn, &tag)
			pid, _ := eu.rr.CreatePath(&path)
			if pid == -1 {
				log.Println("Pid = -1")
				break
			}
			path.Pid = pid

			log.Println("path : ", path)
		}
	} else if mn.Type == "TAG" {
		stations, _ := eu.rr.FindNodesBySinkID(STATION)
		for _, station := range stations {
			log.Println("station : ", station.Name)
			
			path := adapter.PathToModel(&station, &mn)
			pid, _ := eu.rr.CreatePath(&path)
			if pid == -1 {
				log.Println("Pid = -1")
				break
			}
			path.Pid = pid

			log.Println("path : ", path)
		}
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
		log.Println("model.Logic : ", ml)
		return eu.ls.CreateAndStartLogic(&ml)
	}
}

func (eu *eventUsecase) DeleteLogic(l *adapter.Logic) error {
	return eu.ls.RemoveLogic(l.NodeID, l.ID)
}

func (eu *eventUsecase) CreateDelivery(d *adapter.Delivery) error {
	md := adapter.DeliveryToModel(d)
	return eu.rr.CreateDelivery(d.ID, &md)
}
