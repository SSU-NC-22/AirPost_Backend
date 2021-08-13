package registUsecase

import "github.com/eunnseo/AirPost/application/domain/repository"

type registUsecase struct {
	sir repository.SinkRepo
	ndr repository.NodeRepo
	// snr repository.SensorRepo
	lgr repository.LogicRepo
	lsr repository.LogicServiceRepo
	tpr repository.TopicRepo
	acr repository.ActuatorRepo
	dlr repository.DeliveryRepo
}

func NewRegistUsecase(sir repository.SinkRepo,
	ndr repository.NodeRepo,
	// snr repository.SensorRepo,
	lgr repository.LogicRepo,
	lsr repository.LogicServiceRepo,
	tpr repository.TopicRepo,
	acr repository.ActuatorRepo,
	dlr repository.DeliveryRepo) *registUsecase {
	return &registUsecase{
		sir: sir,
		ndr: ndr,
		// snr: snr,
		lgr: lgr,
		lsr: lsr,
		tpr: tpr,
		acr: acr,
		dlr: dlr,
	}
}
