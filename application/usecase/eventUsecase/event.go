package eventUsecase

import (
	"log"
	"sync"

	"github.com/eunnseo/AirPost/application/adapter"
	"github.com/eunnseo/AirPost/application/domain/model"
)

func waitRespGroup(e EVENT, body interface{}, ll []model.LogicService) (prl []pingRequest) {
	var wg sync.WaitGroup // 고루틴이 종료될 때까지 대기
	for _, l := range ll {
		wg.Add(1) // WaitGroup에 대기 중인 고루틴 개수 추가
		go func(_l model.LogicService) {
			url := makeUrl(_l.Addr, EventPath[e]) // 이벤트 발생을 위한 url 생성
			resp, _ := eventClient.R().SetBody(body).Post(url) // POST 수행
			log.Println("Post 내용 : ", body, "url : ", url)
			if !resp.IsSuccess() {
				prl = append(prl, pingRequest{_l, e, body})
			}
			wg.Done() // 대기 중인 고루틴의 수행이 종료되는 것을 알려줌
		}(l)
	}
	wg.Wait() // 모든 고루틴이 종료될 때까지 대기
	return
}

/**************************************************************/
/* sink event usecase                                         */
/**************************************************************/
func (eu *eventUsecase) DeleteSinkEvent(s *model.Sink) error {
	e := DeleteSink

	ll, err := eu.lsr.FindsByTopicID(s.Topic.ID)
	if err != nil {
		return err
	}

	eu.requestRetry = append(eu.requestRetry, waitRespGroup(e, s.Nodes, ll)...)
	// var wg sync.WaitGroup
	// for _, l := range ll {
	// 	wg.Add(1)
	// 	go func() {
	// 		url := makeUrl(l.Addr, path)
	// 		eventClient.R().SetBody(s.Nodes).Post(url)
	// 	}()
	// }
	// wg.Wait()

	return nil
}
func (eu *eventUsecase) CreateSinkEvent(s *model.Sink) error {
	e := CreateSink
	sinkaddr := adapter.SinkAddr{
		Sid:  s.ID,
		Addr: s.Addr,
	}

	ll, err := eu.lsr.FindsByTopicID(s.Topic.ID)
	if err != nil {
		return err
	}
	eu.requestRetry = append(eu.requestRetry, waitRespGroup(e, sinkaddr, ll)...)

	return nil
}

/**************************************************************/
/* node event usecase                                         */
/**************************************************************/
func (eu *eventUsecase) CreateNodeEvent(n *model.Node) error {
	e := CreateNode

	ll, err := eu.lsr.FindsByTopicID(n.Sink.Topic.ID)
	if err != nil {
		return err
	}
	eu.requestRetry = append(eu.requestRetry, waitRespGroup(e, *n, ll)...)

	return nil
}

func (eu *eventUsecase) DeleteNodeEvent(n *model.Node) error {
	e := DeleteNode

	ll, err := eu.lsr.FindsByTopicID(n.Sink.Topic.ID)
	if err != nil {
		return err
	}
	eu.requestRetry = append(eu.requestRetry, waitRespGroup(e, *n, ll)...)

	return nil
}

/**************************************************************/
/* sensor event usecase                                       */
/**************************************************************/
func (eu *eventUsecase) DeleteSensorEvent(s *model.Sensor) error {
	e := DeleteSensor

	ll, err := eu.lsr.Finds()
	if err != nil {
		return err
	}
	eu.requestRetry = append(eu.requestRetry, waitRespGroup(e, *s, ll)...)

	return nil
}

/**************************************************************/
/* logic event usecase                                         */
/**************************************************************/
func (eu *eventUsecase) CreateLogicEvent(l *model.Logic) error {
	e := CreateLogic

	ll, err := eu.lsr.Finds()
	if err != nil {
		return err
	}
	eu.requestRetry = append(eu.requestRetry, waitRespGroup(e, *l, ll)...)

	return nil
}

func (eu *eventUsecase) DeleteLogicEvent(l *model.Logic) error {
	e := DeleteLogic

	ll, err := eu.lsr.Finds()
	if err != nil {
		return err
	}
	eu.requestRetry = append(eu.requestRetry, waitRespGroup(e, *l, ll)...)

	return nil
}
