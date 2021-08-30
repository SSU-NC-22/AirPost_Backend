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
		log.Println("in NewLogicCoreUsecase, run go routin")
		for rawData := range in {
			log.Println("2")

			ld, err := lcu.ToLogicData(&rawData)
			if err != nil {
				log.Println("Error in NewLogicCoreUsecase in ToLogicData")
				continue
			}
			log.Println("in NewLogicCoreUsecase, ld = ", ld)

			lchs, err := lcu.ls.GetLogicChans(ld.NodeID)
			if err == nil {
				log.Println("in NewLogicCoreUsecase, lchs = ", lchs)
				for _, ch := range lchs {
					log.Println("?????")
					if len(ch) != cap(ch) {
						log.Println("?????-----?????")
						ch <- ld // go to "listen" in CreateAndStartLogic core.go
					}
				}
			}
			log.Println("in NewLogicCoreUsecase, after lchs <- ?????")
			out <- lcu.ToDocument(&ld) // go to 
		}
	}()

	return lcu
}

func (lcu *logicCoreUsecase) AppendSinkAddr(sa *adapter.SinkAddr) error {
	lcu.rr.AppendSinkAddr(sa.Sid, &sa.Addr)

	return nil
}
