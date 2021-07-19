package registUsecase

import "github.com/eunnseo/AirPost/application/domain/model"

/**************************************************************/
/* logic service regist usecase                               */
/**************************************************************/
func (ru *registUsecase) GetLogicServices() ([]model.LogicService, error) {
	return ru.lsr.FindsWithTopic()
}

func (ru *registUsecase) UnregistLogicService(l *model.LogicService) error {
	return ru.lsr.Delete(l)
}

/**************************************************************/
/* topic regist usecase                                       */
/**************************************************************/
func (ru *registUsecase) GetTopics() ([]model.Topic, error) {
	return ru.tpr.FindsWithLogicService()
}

func (ru *registUsecase) RegistTopic(t *model.Topic) error {
	return ru.tpr.Create(t)
}

func (ru *registUsecase) UnregistTopic(t *model.Topic) error {
	return ru.tpr.Delete(t)
}
