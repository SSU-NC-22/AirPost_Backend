package adapter

import (
	"math"
	"testing"
)

func TestHaversine(t *testing.T) {
	tests := []struct {
		name                   string
		lat1, lon1, lat2, lon2 float64
		want                   float64 // meters
		tol                    float64 // tolerance in meters
	}{
		{
			name: "same point is zero",
			lat1: 37.5665, lon1: 126.9780,
			lat2: 37.5665, lon2: 126.9780,
			want: 0, tol: 0.001,
		},
		{
			name: "seoul to busan ~325km",
			lat1: 37.5665, lon1: 126.9780,
			lat2: 35.1796, lon2: 129.0756,
			want: 325000, tol: 5000,
		},
		{
			name: "one degree of latitude ~111km",
			lat1: 0, lon1: 0,
			lat2: 1, lon2: 0,
			want: 111195, tol: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Haversine(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.Abs(got-tt.want) > tt.tol {
				t.Errorf("Haversine() = %.2f m, want %.2f m (tol %.2f)", got, tt.want, tt.tol)
			}
		})
	}
}
