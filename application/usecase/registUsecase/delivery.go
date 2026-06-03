package registUsecase

import (
	"errors"
	"math"

	"github.com/eunnseo/AirPost/application/domain/model"
)

// stationSinkID is the sink that classifies station nodes (mirrors handler STATION).
const stationSinkID = 2

/**************************************************************/
/* delivery regist usecase                                    */
/**************************************************************/
func (ru *registUsecase) GetDeliveryByOrderNum(on string) (model.Delivery, error) {
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
// GetShortestPathStation returns the station nearest the drop tag — computed
// GEOMETRICALLY from the registered node coordinates (haversine), so it works for
// any dynamically self-registered topology with NO pre-seeded Path rows. (The old
// implementation required hardcoded tag->station Path distances; that coupling is
// removed so IoT devices that register themselves are immediately routable.)
func (ru *registUsecase) GetShortestPathStation(tagid int) (*model.Node, error) {
	tag, err := ru.ndr.FindsByID(tagid)
	if err != nil {
		return nil, err
	}
	stations, err := ru.ndr.FindsBySinkIDWithSensorValues(stationSinkID)
	if err != nil {
		return nil, err
	}

	var nearest *model.Node
	best := math.MaxFloat64
	for i := range stations {
		s := &stations[i]
		d := Haversine(tag.LocLat, tag.LocLon, s.LocLat, s.LocLon)
		if d < best {
			best, nearest = d, s
		}
	}
	if nearest == nil {
		return nil, errors.New("no station registered to land at")
	}
	return nearest, nil
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
