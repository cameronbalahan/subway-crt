package tui

import (
	"strings"
	"testing"
	"time"

	"subway-crt/internal/transit"
)

func TestRenderDepartureBoardNominal(t *testing.T) {
	styles := MakeStyles(PaletteOutrun)
	snap := transit.TransitSnapshot{
		Timestamp: time.Now(),
		ManhattanTrains: []transit.TrainDeparture{
			{RouteID: "L", Destination: "8 Av", MinutesAway: 2, SecondsAway: 120, CurrentStatus: "Approaching Bedford"},
		},
		CanarsieTrains: []transit.TrainDeparture{
			{RouteID: "L", Destination: "Rockaway Pkwy", MinutesAway: 5, SecondsAway: 300, CurrentStatus: "At 1 Av"},
		},
		ActiveAlerts: []transit.TransitAlert{
			{Severity: "INFO", Summary: "Good Service", AffectsL: true},
		},
		IsStale: false,
	}

	out := RenderDepartureBoard(120, snap, styles)
	if !strings.Contains(out, "MANHATTAN-BOUND") {
		t.Errorf("expected header MANHATTAN-BOUND in output")
	}
	if !strings.Contains(out, "CANARSIE-BOUND") {
		t.Errorf("expected header CANARSIE-BOUND in output")
	}
	if !strings.Contains(out, "2 MIN") {
		t.Errorf("expected '2 MIN' countdown in output")
	}
}

func TestRenderDepartureBoardCriticalAlert(t *testing.T) {
	styles := MakeStyles(PaletteOutrun)
	snap := transit.TransitSnapshot{
		Timestamp: time.Now(),
		ActiveAlerts: []transit.TransitAlert{
			{Severity: "CRITICAL", Summary: "No L Trains Between 8 Av and Lorimer", AffectsL: true},
		},
		IsStale: false,
	}

	out := RenderDepartureBoard(120, snap, styles)
	if !strings.Contains(out, "CRITICAL MTA ADVISORY") {
		t.Errorf("expected critical alert banner in output")
	}
}

func TestRenderDepartureBoardStale(t *testing.T) {
	styles := MakeStyles(PaletteOutrun)
	snap := transit.TransitSnapshot{
		Timestamp:       time.Now(),
		LastSuccessTime: time.Now().Add(-60 * time.Second),
		IsStale:         true,
	}

	out := RenderDepartureBoard(120, snap, styles)
	if !strings.Contains(out, "LIVE FEED DEGRADED") {
		t.Errorf("expected degraded banner when stale")
	}
}
