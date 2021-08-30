package logicCoreUC

import (
	"strings"
	"log"

	"github.com/eunnseo/AirPost/logic-core/domain/model"
)

func (lcuc *logicCoreUsecase) ToLogicData(kd *model.KafkaData) (model.LogicData, error) {
	n, err := lcuc.rr.FindNode(kd.NodeID)
	if err != nil {
		log.Println("Error in ToLogicData from lcuc.rr.FindNode(kd.NodeID)")
		return model.LogicData{}, err
	}
	// s, err := lcuc.rr.FindSensor(kd.SensorID)
	// if err != nil {
	// 	return model.LogicData{}, err
	// }

	vl := map[string]float64{}
	for i, v := range n.SensorValues {
		vl[v] = kd.Values[i]
	}
	return model.LogicData{
		// SensorID:   kd.SensorID,
		// SensorName: s.Name,
		Values:     vl,
		NodeID:		kd.NodeID,
		Node:       *n,
		Timestamp:  kd.Timestamp,
	}, nil
}

func (lcuc *logicCoreUsecase) ToDocument(ld *model.LogicData) model.Document {
	sinkname := ld.Node.SinkName
	if sinkname[0]==' '{
		sinkname=sinkname[1:]
	}
	return model.Document{		
		Index: "airpost-" + (strings.Split(ld.Node.Name,"-"))[1]+"-" + strings.ReplaceAll(sinkname," ", "-"),
		Doc:   *ld,
	}
}