package repository

import (
	"github.com/eunnseo/AirPost/logic-core/domain/model"
)

type RegistRepo interface {
	FindNode(key int) (*model.Node, error)
	FindNodesBySinkID(sid int) ([]model.Node, error)
	CreateNode(key int, n *model.Node) error
	DeleteNode(key int) error
	
	AppendSinkAddr(sid int, s *string) error

	FindPath(key int) (*model.Path, error)
	CreatePath(p *model.Path) (int, error)
	DeletePath(key int) error
	FindShortestPathStation(tagid int) (station *model.Node, err error)

	FindDelivery(key int) (*model.Delivery, error)
	CreateDelivery(key int, d *model.Delivery) error
	DeleteDelivery(key int) error
}
