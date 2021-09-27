package service

import "github.com/eunnseo/AirPost/logic-core/domain/model"

type KafkaConsumerGroup interface {
	GetOutput() <-chan model.KafkaData
}
