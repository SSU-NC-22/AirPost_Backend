package registUsecase

import (
	"testing"

	"github.com/eunnseo/AirPost/application/domain/model"
	"github.com/eunnseo/AirPost/application/domain/repository"
)

// fakeNodeRepo serves a drop tag (FindsByID) and the registered stations
// (FindsBySinkIDWithSensorValues) so GetShortestPathStation can be driven purely
// from node coordinates — no seeded Path rows.
type fakeNodeRepo struct {
	repository.NodeRepo
	tag      *model.Node
	stations []model.Node
}

func (f *fakeNodeRepo) FindsByID(id int) (*model.Node, error) {
	if f.tag != nil && f.tag.ID == id {
		return f.tag, nil
	}
	return &model.Node{ID: id}, nil
}

func (f *fakeNodeRepo) FindsBySinkIDWithSensorValues(sinkid int) ([]model.Node, error) {
	return f.stations, nil
}

var _ repository.NodeRepo = (*fakeNodeRepo)(nil)

func TestGetShortestPathStation(t *testing.T) {
	// A drop tag near (37.5000, 127.0000).
	tag := &model.Node{ID: 1, LocLat: 37.5000, LocLon: 127.0000}

	tests := []struct {
		name       string
		stations   []model.Node
		wantErr    bool
		wantNodeID int
	}{
		{
			name:     "no registered stations returns error",
			stations: nil,
			wantErr:  true,
		},
		{
			name: "picks the geometrically nearest station",
			stations: []model.Node{
				{ID: 10, LocLat: 37.5100, LocLon: 127.0000}, // ~1.1 km north
				{ID: 20, LocLat: 37.5005, LocLon: 127.0000}, // ~55 m north (nearest)
				{ID: 30, LocLat: 37.4000, LocLon: 127.5000}, // far
			},
			wantNodeID: 20,
		},
		{
			name:       "single station",
			stations:   []model.Node{{ID: 7, LocLat: 37.6, LocLon: 127.1}},
			wantNodeID: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ru := &registUsecase{ndr: &fakeNodeRepo{tag: tag, stations: tt.stations}}
			station, err := ru.GetShortestPathStation(tag.ID)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (station=%v)", station)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if station == nil || station.ID != tt.wantNodeID {
				t.Errorf("got station %v, want ID %d", station, tt.wantNodeID)
			}
		})
	}
}
