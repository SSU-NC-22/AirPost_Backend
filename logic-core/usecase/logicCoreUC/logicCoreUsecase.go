package logicCoreUC

import (
	"github.com/eunnseo/AirPost/logic-core/adapter"
	"github.com/eunnseo/AirPost/logic-core/domain/repository"
	"github.com/eunnseo/AirPost/logic-core/domain/service"
)

type logicCoreUsecase struct {
	rr repository.RegistRepo
	ks service.KafkaConsumerGroup
	es service.ElasticClient
	ls service.LogicService
}

func NewLogicCoreUsecase(rr repository.RegistRepo,
	ks service.KafkaConsumerGroup,
	es service.ElasticClient,
	ls service.LogicService) *logicCoreUsecase {
	lcu := &logicCoreUsecase{
		rr: rr,
		ks: ks,
		es: es,
		ls: ls,
	}

	in := lcu.ks.GetOutput()
	out := lcu.es.GetInput()

	go func() {
		for rawData := range in {

			ld, err := lcu.ToLogicData(&rawData) // 데이터 보강
			if err != nil {
				continue // unknown node or value/schema mismatch
			}

			// Route to any active logic chains for this node (delivery tracking, arrival
			// alarms). Plain telemetry has no logic attached — that is normal, not an error.
			if lchs, err := lcu.ls.GetLogicChans(ld.NodeID); err == nil {
				for _, ch := range lchs {
					if len(ch) != cap(ch) {
						ch <- ld // go to "listen" in CreateAndStartLogic core.go
					}
				}
			}

			out <- lcu.ToDocument(&ld) // always archive the reading to Elasticsearch
		}
	}()

	return lcu
}

func (lcu *logicCoreUsecase) AppendSinkAddr(sa *adapter.SinkAddr) error {
	lcu.rr.AppendSinkAddr(sa.Sid, &sa.Addr)

	return nil
}
