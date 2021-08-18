package service

import "github.com/eunnseo/AirPost/logic-core/domain/model"

type LogicService interface {
	CreateAndStartLogic(l *model.Logic) error
	RemoveLogic(nid, lid int) error
	GetLogicChans(nid int) (map[int]chan model.LogicData, error)
	
}
