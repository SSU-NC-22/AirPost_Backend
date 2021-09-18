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

/*
func (eu *eventUsecase) CreateDeliveryLogic(d *adapter.Delivery) error {
	log.Println("===== eu.CreateDeliveryLogic start =====")

	srcStation, err := eu.rr.FindNode(d.SrcStationID)
	if err != nil {
		log.Println("no src station")
		return err
	}
	log.Println("srcStation : ", srcStation)

	destTag, err := eu.rr.FindNode(d.DestStationID)
	if err != nil {
		log.Println("no dest tag")
		return err
	}
	log.Println("destTag : ", destTag)
	
	destStationID, err := eu.rr.FindShortestPathStationID(destTag.Nid)
	if err != nil {
		log.Println("no dest station")
		return err
	}
	destStation, err := eu.rr.FindNode(destStationID)
	if err != nil {
		log.Println("no dest station")
		return err
	}
	log.Println("destStation : ", destStation)

	srcStationLoc := []float64{37.497365670723944, 126.95591909983503} // 위도(lat), 경도(lon)
	tagLoc := []float64{37.49719755738831, 126.95590032437323}
	destStationLoc := []float64{37.4971933013496, 126.95619804955307}

	// srcStationLoc := []float64{srcStation.Location.Lat, srcStation.Location.Lon} // 위도(lat), 경도(lon)
	// tagLoc := []float64{destTag.Location.Lat, destTag.Location.Lon}
	// destStationLoc := []float64{destStation.Location.Lat, destStation.Location.Lon}

	// path 초기화
	path := [][]float64{}
	path = append(path, srcStationLoc)
	path = append(path, tagLoc)
	path = append(path, destStationLoc)

	// drone event element 생성
	me := model.Element{
		Elem: "drone",
		Arg:  map[string]interface{} {
			"nid": "DRO" + strconv.Itoa(d.DroneID),
			"values": path,
			"tagidx": 1, // TODO
		},
	}

	// drone event logic 생성
	ml := model.Logic{
		LogicName: "drone",
		Elems:	   []model.Element{me},
		NodeID:	   d.DroneID,
	}
	adapter.LogicToModel(&ml)

	log.Println("ml : ", ml)
	return eu.ls.CreateAndStartLogic(&ml)
}
*/

func (eu *eventUsecase) CreateDelivery(d *adapter.Delivery) error {
	md := adapter.DeliveryToModel(d)
	return eu.rr.CreateDelivery(d.ID, &md)
}
