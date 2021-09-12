package registUsecase

import (
	"github.com/eunnseo/AirPost/application/domain/model"
)

/**************************************************************/
/* delivery regist usecase                                    */
/**************************************************************/
func (ru *registUsecase) GetDeliveryByOrderNum(on int) (model.Delivery, error) {
	return ru.dlr.FindsByOrderNum(on)
}

func (ru *registUsecase) RegistDelivery(d *model.Delivery) error {
	return ru.dlr.Create(d)
}

func (ru *registUsecase) UnregistDelivery(d *model.Delivery) error {
	return ru.dlr.Delete(d)
}

/**************************************************************/
/* path regist usecase                                        */
/**************************************************************/
func (ru *registUsecase) RegistPath(p *model.Path) error {
	return ru.ptr.Create(p)
}

func (ru *registUsecase) UnregistPath(p *model.Path) error {
	return ru.ptr.Delete(p)
}

/**************************************************************/
/* drone_loc regist usecase                                   */
/**************************************************************/
func (ru *registUsecase) GetStationDroneByStationID(stationid int) ([]model.StationDrone, error) {
	return ru.sdr.FindsByStationID(stationid)
}

func (ru *registUsecase) RegistStationDrone(sd *model.StationDrone) error {
	sd.Reserved = 1
	return ru.sdr.Create(sd)
}

func (ru *registUsecase) UnregistStationDrone(sd *model.StationDrone) error {
	sd.Reserved = 0
	return ru.sdr.Delete(sd)
}
