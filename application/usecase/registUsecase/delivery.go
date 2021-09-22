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
func (ru *registUsecase) GetShortestPathStation(tagid int) (station *model.Node, err error) {
	pl, err := ru.ptr.Finds()
	if err != nil {
		return nil, err
	}
	// to do : if pl empty
	min := pl[0].Distance
	nid := pl[0].StationID
	for _, path := range(pl) {
		if (path.Distance < min) {
			min = path.Distance
			nid = path.StationID
		}
	}
	station, err = ru.ndr.FindsByID(nid)
	if err != nil {
		return nil, err
	}
	return station, nil
}

func (ru *registUsecase) GetPaths() ([]model.Path, error) {
	return ru.ptr.Finds()
}

func (ru *registUsecase) RegistPath(p *model.Path) error {
	return ru.ptr.Create(p)
}

func (ru *registUsecase) UnregistPath(p *model.Path) error {
	return ru.ptr.Delete(p)
}

/**************************************************************/
/* station_drone regist usecase                               */
/**************************************************************/
func (ru *registUsecase) GetStationDrone(stationid int, droneid int) (*model.StationDrone, error) {
	return ru.sdr.Find(stationid, droneid)
}

func (ru *registUsecase) GetStationDroneByStationID(stationid int) ([]model.StationDrone, error) {
	return ru.sdr.FindsByStationID(stationid)
}

func (ru *registUsecase) GetStationDroneByDroneID(droneid int) ([]model.StationDrone, error) {
	return ru.sdr.FindsByDroneID(droneid)
}

func (ru *registUsecase) RegistStationDrone(sd *model.StationDrone) error {
	return ru.sdr.Create(sd)
}

func (ru *registUsecase) UnregistStationDrone(sd *model.StationDrone) error {
	return ru.sdr.Delete(sd)
}

func (ru *registUsecase) UnregistStationDroneByStationID(sd *model.StationDrone) error {
	return ru.sdr.DeleteByStationID(sd)
}

func (ru *registUsecase) UnregistStationDroneByDroneID(sd *model.StationDrone) error {
	return ru.sdr.DeleteByDroneID(sd)
}
