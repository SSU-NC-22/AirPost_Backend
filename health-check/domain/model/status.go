package model

import (
	"log"
	"time"

	"github.com/eunnseo/AirPost/health-check/adapter"
	// "github.com/eunnseo/AirPost/health-check/setting"
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
	Battery  int       `json:"battery"`
	Location []float64 `json:"location"` // [lat, lon, alt]
}

type Status struct {
	State       int       `json:"state"`
	Work        bool      `json:"work"`
	Battery     int       `json:"battery"`
	Location    []float64 `json:"location"`
	LastConnect time.Time `json:"last_connect"`
}

func NewStatus(ns adapter.NodeState, t time.Time) Status { // 인자로 받은 work 여부로 Status 구조체 설정
	res := Status{
		Work:        ns.State,
		Battery:     ns.Battery,
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

func (s *Status) setState(v int) { // 인자로 받은 v로 Status구조체 변경
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
	if s.Battery != ns.Battery {
		log.Println("changed battery")
		s.Battery = ns.Battery
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
	timeout := time.Now() //s.LastConnect.Add(time.Duration(setting.StatusSetting.Drop) * time.Hour)
	return now.After(timeout) // TODO
}
