package registUsecase

import (
	"testing"

	"github.com/eunnseo/AirPost/application/domain/model"
	"github.com/eunnseo/AirPost/application/domain/repository"
)

// fakePathRepo lets us drive GetShortestPathStation with a fixed path set.
type fakePathRepo struct {
	repository.PathRepo
	paths []model.Path
	err   error
}

func (f *fakePathRepo) Finds() ([]model.Path, error) { return f.paths, f.err }

// fakeNodeRepo returns a fixed node for FindsByID.
type fakeNodeRepo struct {
	repository.NodeRepo
}

func (f *fakeNodeRepo) FindsByID(id int) (*model.Node, error) {
	return &model.Node{ID: id}, nil
}

// compile-time guards that our fakes satisfy the interfaces.
var (
	_ repository.PathRepo = (*fakePathRepo)(nil)
	_ repository.NodeRepo = (*fakeNodeRepo)(nil)
)

func TestGetShortestPathStation(t *testing.T) {
	tests := []struct {
		name       string
		paths      []model.Path
		tagid      int
		wantErr    bool
		wantNodeID int
	}{
		{
			name:    "empty path slice does not panic, returns error",
			paths:   nil,
			tagid:   1,
			wantErr: true,
		},
		{
			name: "no path matches tag id returns error",
			paths: []model.Path{
				{StationID: 10, TagID: 99, Distance: 5},
			},
			tagid:   1,
			wantErr: true,
		},
		{
			name: "picks nearest station for the tag",
			paths: []model.Path{
				{StationID: 10, TagID: 1, Distance: 50},
				{StationID: 20, TagID: 1, Distance: 5},
				{StationID: 30, TagID: 2, Distance: 1}, // different tag, ignored
			},
			tagid:      1,
			wantErr:    false,
			wantNodeID: 20,
		},
		{
			name: "single matching path",
			paths: []model.Path{
				{StationID: 7, TagID: 3, Distance: 12},
			},
			tagid:      3,
			wantErr:    false,
			wantNodeID: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ru := &registUsecase{
				ptr: &fakePathRepo{paths: tt.paths},
				ndr: &fakeNodeRepo{},
			}
			station, err := ru.GetShortestPathStation(tt.tagid)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (station=%v)", station)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if station == nil {
				t.Fatalf("expected station, got nil")
			}
			if station.ID != tt.wantNodeID {
				t.Errorf("got station ID %d, want %d", station.ID, tt.wantNodeID)
			}
		})
	}
}
