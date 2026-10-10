package transit

import (
	"testing"
	"time"
)

func TestFallbackSimulation(t *testing.T) {
	now := time.Now()
	snap := FallbackSimulation(now, false)

	if len(snap.ManhattanTrains) == 0 {
		t.Fatalf("expected Manhattan trains, got 0")
	}
	if len(snap.CanarsieTrains) == 0 {
		t.Fatalf("expected Canarsie trains, got 0")
	}
	if snap.IsStale {
		t.Fatalf("expected not stale")
	}

	for _, dep := range snap.ManhattanTrains {
		if dep.Destination == "" {
			t.Errorf("expected destination to be populated")
		}
		if dep.SecondsAway < 0 {
			t.Errorf("expected non-negative seconds away, got %d", dep.SecondsAway)
		}
	}
}

func TestMapStationToGrid(t *testing.T) {
	tests := []struct {
		stopID string
		wantX  int
		wantY  int
	}{
		{"L01", 15, 18},
		{"L10", 126, 18},
		{"L11", 146, 20},
		{"G29", 146, 20},
	}

	for _, tc := range tests {
		col, row := MapStationToGrid(tc.stopID, "L", "N", "STOPPED_AT")
		if col != tc.wantX || row != tc.wantY {
			t.Errorf("stop %s: expected (%d, %d), got (%d, %d)", tc.stopID, tc.wantX, tc.wantY, col, row)
		}
	}
}
