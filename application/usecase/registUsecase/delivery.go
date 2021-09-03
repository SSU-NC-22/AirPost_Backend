package registUsecase

import (
	"github.com/eunnseo/AirPost/application/domain/model"
)

/**************************************************************/
/* delivery regist usecase                                    */
/**************************************************************/
func (ru *registUsecase) RegistDelivery(d *model.Delivery) error {
	return ru.dlr.Create(d)
}

// func (ru *registUsecase) UnregistDelivery(d *model.Delivery) error {
// 	return ru.dlr.Delete(d)
// }

func (ru *registUsecase) GetDeliveryByOrderNum(on int) (model.Delivery, error) {
	return ru.dlr.FindsByOrderNum(on)
}
