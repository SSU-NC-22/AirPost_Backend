package model

import (
	"log"
	"time"
	"math"

	"github.com/eunnseo/AirPost/health-check/adapter"
	"github.com/eunnseo/AirPost/health-check/setting"
)

const (
	RED    = 0 // 미동작
	YELLOW = 1 // 형태 바뀔 경우 가운데서 중재 단계
	GREEN  = 2 // 동작
)

type SinkStatus struct {
	SinkID  int          `json:"sid"`
	Satates []NodeStatus `json:"states"`
}

type NodeStatus struct {
	NodeID   int       `json:"nid"`
	State    int       `json:"state"`
	Battery  int       `json:"battery"` // 퍼센트
	Location []float64 `json:"location"` // [lat, lon, alt]
}

type Status struct {
	State       int       `json:"state"`
	Work        bool      `json:"work"`
	Battery     int       `json:"battery"`
	Location    []float64 `json:"location"`
	LastConnect time.Time `json:"last_connect"`
}

func ToPercentage(battery float64) int {
	minVal := 14.0
	maxVal := 16.8
	per := math.Round((battery - minVal) / (maxVal-minVal) * 100)
	return int(per)
}

func NewStatus(ns adapter.NodeState, t time.Time) Status {
	res := Status{
		Work:        ns.State,
		Battery:     ToPercentage(ns.Battery),
		Location:    ns.Location,
		LastConnect: t,
	}
	if ns.State {
		res.State = GREEN
	} else {
		res.State = RED
	}
	return res
}

func (s *Status) setState(v int) {
	s.State = v
	switch v {
	case RED:
		s.Work = false
	case GREEN:
		s.Work = true
	case YELLOW:
		s.Work = !s.Work
	}
}

func (s *Status) UpdateState(ns adapter.NodeState, t time.Time) bool {
	isChanged := false

	// Update time for drop
	if ns.State {
		s.LastConnect = t
	}
	if s.State == YELLOW {
		if ns.State {
			s.setState(GREEN)
		} else {
			s.setState(RED)
		}
		isChanged = true
	} else if s.Work != ns.State {
		s.setState(YELLOW)
		isChanged = true
	}

	// Update battery
	battery := ToPercentage(ns.Battery)
	if s.Battery != battery {
		log.Println("changed battery")
		s.Battery = battery
		isChanged = true
	}

	// Update location
	for i, loc := range ns.Location {
		if s.Location[i] != loc {
			log.Println("changed location")
			s.Location[i] = loc
			isChanged = true
		}
	}

	return isChanged
}

func (s *Status) CheckDrop() bool {
	s.setState(RED)
	now := time.Now()
	timeout := s.LastConnect.Add(time.Duration(setting.StatusSetting.Drop) * time.Hour)
	return now.After(timeout)
}
