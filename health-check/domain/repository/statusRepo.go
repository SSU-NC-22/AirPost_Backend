package repository

import (
	"github.com/eunnseo/AirPost/health-check/adapter"
	"github.com/eunnseo/AirPost/health-check/domain/model"
)

type StatusRepo interface {
	UpdateTable(states adapter.States) []model.NodeStatus
}
