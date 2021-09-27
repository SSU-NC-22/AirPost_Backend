package logicCoreUC

import (
	"log"

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
				continue
			}

			lchs, err := lcu.ls.GetLogicChans(ld.NodeID)
			if err != nil {
				if ld.Node.Type == "DRO" {
					log.Println("it's drone") // delivery 없음
					continue
				} else {
					panic(err)
				}
			}
			if err == nil {
				for _, ch := range lchs {
					if len(ch) != cap(ch) {
						ch <- ld // go to "listen" in CreateAndStartLogic core.go
					}
				}
			}
			out <- lcu.ToDocument(&ld) // go to elastic client
		}
	}()

	return lcu
}

func (lcu *logicCoreUsecase) AppendSinkAddr(sa *adapter.SinkAddr) error {
	lcu.rr.AppendSinkAddr(sa.Sid, &sa.Addr)

	return nil
}
