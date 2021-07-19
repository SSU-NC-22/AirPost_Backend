package registUsecase

import "github.com/eunnseo/AirPost/application/domain/model"

/**************************************************************/
/* logic regist usecase                                       */
/**************************************************************/
func (ru *registUsecase) GetLogics() ([]model.Logic, error) {
	return ru.lgr.FindsWithSensorValues()
}

func (ru *registUsecase) RegistLogic(l *model.Logic) error {
	return ru.lgr.Create(l)
}

func (ru *registUsecase) UnregistLogic(l *model.Logic) error {
	return ru.lgr.Delete(l)
}
