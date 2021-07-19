package service

import "github.com/eunnseo/AirPost/logic-core/domain/model"

type ElasticClient interface {
	GetInput() chan<- model.Document
}
