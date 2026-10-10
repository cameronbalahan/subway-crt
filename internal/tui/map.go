package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"subway-crt/internal/transit"
)

// RenderSchematicMap renders the bottom schematic transit map with Outrun radar aesthetics
func RenderSchematicMap(width, height int, snap transit.TransitSnapshot, styles Styles) string {
	var b strings.Builder

	titleBar := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(styles.Palette.Secondary).
		Padding(0, 2).
		Render("▲ TRANSIT SCHEMATIC RADAR // MANHATTAN & NORTH BROOKLYN SUBWAY CORRIDOR ▲")

	b.WriteString(titleBar)
	b.WriteString("\n")

	// Map real train positions to quick lookup indicators
	beaconAt := make(map[string]string)
	for _, v := range snap.TrainPositions {
		cleanID := v.StopID
		if len(cleanID) > 3 {
			cleanID = cleanID[:3]
		}
		glyph := "●"
		if v.Direction == "N" {
			glyph = "▲"
		} else {
			glyph = "▼"
		}
		beaconAt[cleanID] = glyph
	}

	getBeacon := func(stopID string) string {
		if g, ok := beaconAt[stopID]; ok {
			return styles.FgAccent.Render(g)
		}
		return styles.FgMuted.Render("○")
	}

	// Geographic Header with retro brackets
	headerGeo := fmt.Sprintf(
		"  %s               %s         %s",
		styles.FgSecondary.Render("◄◄ MANHATTAN (14TH ST TRUNK) ◄◄"),
		styles.FgRiver.Render("≈≈ [ EAST RIVER ] ≈≈"),
		styles.FgPrimary.Render("►► NORTH BROOKLYN (WILLIAMSBURG / BUSHWICK) ►►"),
	)
	b.WriteString(headerGeo + "\n\n")

	// Intersecting North Bound Trunks with Authentic MTA Colors
	b.WriteString(fmt.Sprintf("  ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐     ║                 ┌─────────┐   ┌─────────┐   ┌─────────┐\n"))
	b.WriteString(fmt.Sprintf("  │ %s │   │ %s │   │ %s │   │ %s │     ║                 │ %s │   │(CROSSTN)│   │ %s │\n",
		styles.BadgeACE.Render("A·C·E"),
		styles.Badge123.Render("1·2·3"),
		styles.BadgeBDFM.Render("B·D·F"),
		styles.BadgeNQRW.Render("N·Q·R"),
		styles.BadgeG.Render("  G  "),
		styles.BadgeJMZ.Render("J·M·Z"),
	))
	b.WriteString(fmt.Sprintf("  │  8th Av │   │  7th Av │   │  6th Av │   │Union Sq │     ║                 │LorimerSt│   │ Graham  │   │Marcy Av │\n"))
	b.WriteString(fmt.Sprintf("  └────┬────┘   └────┬────┘   └────┬────┘   └────┬────┘     ║                 └────┬────┘   └────┬────┘   └────┬────┘\n"))
	b.WriteString(fmt.Sprintf("       │             │             │             │          ║                      │             │             │\n"))

	// L Line Vector with glowing pips & East River crossing
	st8Av := fmt.Sprintf("[%s 8 Av]", getBeacon("L01"))
	st6Av := fmt.Sprintf("[%s 6 Av]", getBeacon("L02"))
	stUnion := fmt.Sprintf("[%s UnionSq]", getBeacon("L03"))
	st3Av := fmt.Sprintf("[%s 3 Av]", getBeacon("L05"))
	st1Av := fmt.Sprintf("[%s 1 Av]", getBeacon("L06"))

	bedfordBox := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(styles.Palette.Accent).
		Padding(0, 1).
		Render(fmt.Sprintf("★ %s BEDFORD AV ★", getBeacon("L10")))

	stLorimer := fmt.Sprintf("[%s Lorimer]", getBeacon("L11"))
	stGraham := fmt.Sprintf("[%s Graham]", getBeacon("L12"))
	stGrand := fmt.Sprintf("[%s Grand]", getBeacon("L13"))
	stMontrose := fmt.Sprintf("[%s Montrose]", getBeacon("L14"))

	riverTube := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Palette.River).
		Render("≈≈[ 14TH ST TUBE ]≈≈")

	track := fmt.Sprintf("  %s═══%s═══%s═══%s═══%s %s %s═══%s═══%s═══%s═══%s══► (Canarsie)",
		styles.StationPip.Render(st8Av),
		styles.StationPip.Render(st6Av),
		styles.StationPip.Render(stUnion),
		styles.StationPip.Render(st3Av),
		styles.StationPip.Render(st1Av),
		riverTube,
		bedfordBox,
		styles.StationPip.Render(stLorimer),
		styles.StationPip.Render(stGraham),
		styles.StationPip.Render(stGrand),
		styles.StationPip.Render(stMontrose),
	)
	b.WriteString(track + "\n")

	// Intersecting South Bound Trunks
	b.WriteString(fmt.Sprintf("       │             │             │             │          ║                      │             │             │\n"))
	b.WriteString(fmt.Sprintf("  ┌────┴────┐   ┌────┴────┐   ┌────┴────┐   ┌────┴────┐     ║                 ┌────┴────┐   ┌────┴────┐   ┌────┴────┐\n"))
	b.WriteString(fmt.Sprintf("  │ %s │   │ (PATH)  │   │ (M14-S) │   │ (M14-A) │     ║                 │(B'WAY-G)│   │ (B24)   │   │ (B39)   │\n",
		styles.Badge456.Render("4·5·6"),
	))
	b.WriteString(fmt.Sprintf("  │Union Sq │   │  14 St  │   │Cross-Twn│   │East-Side│     ║                 │Metropol.│   │Bus Xfer │   │Williamsb│\n"))
	b.WriteString(fmt.Sprintf("  └─────────┘   └─────────┘   └─────────┘   └─────────┘     ║                 └─────────┘   └─────────┘   └─────────┘\n"))

	// Real Beacon Telemetry Bar
	realCount := len(snap.TrainPositions)
	b.WriteString("\n")
	radarTelemetry := fmt.Sprintf(
		"  %s %s %s",
		styles.FgSecondary.Render("⚡ RADAR TELEMETRY:"),
		styles.FgAccent.Render(fmt.Sprintf("%d LIVE TRAINS DETECTED IN CORRIDOR", realCount)),
		styles.FgMuted.Render("·  ▲ North/Westbound (8 Av)  ·  ▼ South/Eastbound (Canarsie)  ·  ★ Station Kiosk"),
	)
	b.WriteString(radarTelemetry + "\n")

	// Outrun Synthwave Horizon perspective wireframe grid
	grid1 := styles.FgBorder.Render("  / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / /")
	grid2 := styles.FgPrimary.Render("  ═════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════")
	b.WriteString(grid1 + "\n")
	b.WriteString(grid2 + "\n")

	return b.String()
}
