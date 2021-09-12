package registUsecase

import "github.com/eunnseo/AirPost/application/domain/repository"

type registUsecase struct {
	sir repository.SinkRepo
	ndr repository.NodeRepo
	lgr repository.LogicRepo
	lsr repository.LogicServiceRepo
	tpr repository.TopicRepo
	acr repository.ActuatorRepo
	dlr repository.DeliveryRepo
	ptr repository.PathRepo
	sdr repository.StationDroneRepo
}

func NewRegistUsecase(sir repository.SinkRepo,
	ndr repository.NodeRepo,
	lgr repository.LogicRepo,
	lsr repository.LogicServiceRepo,
	tpr repository.TopicRepo,
	acr repository.ActuatorRepo,
	dlr repository.DeliveryRepo,
	ptr repository.PathRepo,
	sdr repository.StationDroneRepo) *registUsecase {
	return &registUsecase{
		sir: sir,
		ndr: ndr,
		lgr: lgr,
		lsr: lsr,
		tpr: tpr,
		acr: acr,
		dlr: dlr,
		ptr: ptr,
		sdr: sdr,
	}
}
