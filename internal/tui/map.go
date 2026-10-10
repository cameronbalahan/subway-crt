package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"subway-crt/internal/transit"
)

// RenderSchematicMap renders the bottom schematic transit map with Amber CRT radar aesthetics
// Fully widened for 1080p terminal displays (180-240 columns) with vertically centered alignment
func RenderSchematicMap(width, height int, snap transit.TransitSnapshot, styles Styles) string {
	var b strings.Builder

	titleBar := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Palette.Background).
		Background(styles.Palette.Primary).
		Padding(0, 2).
		Render("▲ TRANSIT SCHEMATIC RADAR // MANHATTAN & NORTH BROOKLYN SUBWAY CORRIDOR (L LINE TRUNK) ▲")

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

	// 10 Station Columns along the L corridor + 1 East River crossing
	// Width of each trunk box: exactly 15 characters
	// Box format:
	//   ┌─────────────┐
	//   │   (A·C·E)   │
	//   │   8th Av    │
	//   └──────┬──────┘
	//          │
	//     [○ 8th Av]
	//          │
	//   ┌──────┴──────┐
	//   │   (M14-A)   │
	//   │  Crosstown  │
	//   └─────────────┘
	//
	// In amber monochrome: all trunks are themed in crisp amber CRT phosphor

	b.WriteString("\n")

	// Geographic Sector Header (spaced out across screen)
	geoManhattan := styles.FgPrimary.Render("◄◄ MANHATTAN // 14TH STREET TRUNK ◄◄")
	geoRiver := styles.FgBorder.Render("≈≈≈ [ EAST RIVER TUBE ] ≈≈≈")
	geoBrooklyn := styles.FgSecondary.Render("►► NORTH BROOKLYN // WILLIAMSBURG & BUSHWICK TRUNK ►►")

	b.WriteString(fmt.Sprintf("    %-68s %-32s %s\n\n", geoManhattan, geoRiver, geoBrooklyn))

	// Row 1: Top Trunk Box Top Borders (7 across)
	b.WriteString("    ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐     ║║║║     ┌─────────────┐   ┌─────────────┐   ┌─────────────┐\n")

	// Row 2: Top Trunk Badges
	b.WriteString(fmt.Sprintf("    │   %s   │   │   %s   │   │   %s   │   │   %s   │     ║║║║     │   %s   │   │   %s   │   │   %s   │\n",
		styles.BadgeAmberMain.Render("(A·C·E)"),
		styles.BadgeAmberMain.Render("(1·2·3)"),
		styles.BadgeAmberMain.Render("(B·D·F)"),
		styles.BadgeAmberMain.Render("(N·Q·R)"),
		styles.BadgeAmberMain.Render("(  G  )"),
		styles.BadgeAmberDim.Render("(CROSSTN)"),
		styles.BadgeAmberMain.Render("(J·M·Z)"),
	))

	// Row 3: Top Trunk Station Labels
	b.WriteString("    │   8th Av    │   │   7th Av    │   │   6th Av    │   │  Union Sq   │     ║║║║     │  Lorimer St │   │  Graham Av  │   │  Marcy Av   │\n")

	// Row 4: Top Trunk Bottom Borders with connecting stems centered at column 8 of each 15-char box
	b.WriteString("    └──────┬──────┘   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘     ║║║║     └──────┬──────┘   └──────┬──────┘   └──────┬──────┘\n")

	// Row 5: Stem connecting lines directly above station pips
	b.WriteString("           │                 │                 │                 │          ║║║║            │                 │                 │     \n")

	// Row 6: Main L Line Track with centered station pips
	// Each station pip is centered with the 15-char box above it (4 leading spaces + 15 width)
	pip8Av := fmt.Sprintf("  [%s  8 Av]  ", getBeacon("L01"))       // 15 chars
	pip6Av := fmt.Sprintf("  [%s  6 Av]  ", getBeacon("L02"))       // 15 chars
	pipUnion := fmt.Sprintf(" [%s UnionSq] ", getBeacon("L03"))     // 15 chars
	pip3Av := fmt.Sprintf("  [%s  3 Av]  ", getBeacon("L05"))       // 15 chars
	pip1Av := fmt.Sprintf("  [%s  1 Av]  ", getBeacon("L06"))       // 15 chars

	riverTube := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Palette.Border).
		Render("≈≈[ 14 ST TUBE ]≈≈") // 18 chars

	pipBedford := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Palette.Background).
		Background(styles.Palette.Primary).
		Padding(0, 1).
		Render(fmt.Sprintf("★ %s BEDFORD AV ★", getBeacon("L10"))) // 19 chars

	pipLorimer := fmt.Sprintf(" [%s Lorimer] ", getBeacon("L11"))   // 15 chars
	pipGraham := fmt.Sprintf(" [%s Graham]  ", getBeacon("L12"))   // 15 chars
	pipGrand := fmt.Sprintf("  [%s Grand]   ", getBeacon("L13"))   // 15 chars
	pipMontrose := fmt.Sprintf(" [%s Montrose] ", getBeacon("L14")) // 15 chars

	// Track rendered across with perfect amber linkage
	b.WriteString(fmt.Sprintf("    %s═══%s═══%s═══%s═══%s %s %s═══%s═══%s═══%s═══%s══► (Canarsie)\n",
		styles.StationPip.Render(pip8Av),
		styles.StationPip.Render(pip6Av),
		styles.StationPip.Render(pipUnion),
		styles.StationPip.Render(pip3Av),
		styles.StationPip.Render(pip1Av),
		riverTube,
		pipBedford,
		styles.StationPip.Render(pipLorimer),
		styles.StationPip.Render(pipGraham),
		styles.StationPip.Render(pipGrand),
		styles.StationPip.Render(pipMontrose),
	))

	// Row 7: Stem connecting lines directly below station pips
	b.WriteString("           │                 │                 │                 │          ║║║║            │                 │                 │     \n")

	// Row 8: Bottom Trunk Top Borders with connecting stems
	b.WriteString("    ┌──────┴──────┐   ┌──────┴──────┐   ┌──────┴──────┐   ┌──────┴──────┐     ║║║║     ┌──────┴──────┐   ┌──────┴──────┐   ┌──────┴──────┐\n")

	// Row 9: Bottom Trunk Badges
	b.WriteString(fmt.Sprintf("    │   %s   │   │   %s   │   │   %s   │   │   %s   │     ║║║║     │   %s   │   │   %s   │   │   %s   │\n",
		styles.BadgeAmberMain.Render("(4·5·6)"),
		styles.BadgeAmberDim.Render("( PATH )"),
		styles.BadgeAmberDim.Render("(M14-SBS)"),
		styles.BadgeAmberDim.Render("(M14-A)"),
		styles.BadgeAmberMain.Render("(B'WAY-G)"),
		styles.BadgeAmberDim.Render("( B24 )"),
		styles.BadgeAmberDim.Render("( B39 )"),
	))

	// Row 10: Bottom Trunk Station / Bus Transfer Names
	b.WriteString("    │  Union Sq   │   │  14th St    │   │  Cross-Town │   │  East Side  │     ║║║║     │ Metropolitan│   │  Bus Xfer   │   │ Williamsburg│\n")

	// Row 11: Bottom Trunk Box Bottom Borders
	b.WriteString("    └─────────────┘   └─────────────┘   └─────────────┘   └─────────────┘     ║║║║     └─────────────┘   └─────────────┘   └─────────────┘\n")

	// Real Beacon Telemetry Bar
	realCount := len(snap.TrainPositions)
	b.WriteString("\n")
	radarTelemetry := fmt.Sprintf(
		"    %s %s %s",
		styles.FgPrimary.Render("⚡ RADAR TELEMETRY:"),
		styles.FgAccent.Render(fmt.Sprintf("%d LIVE TRAINS DETECTED IN CORRIDOR", realCount)),
		styles.FgMuted.Render("·  ▲ Northbound (8 Av)  ·  ▼ Southbound (Canarsie)  ·  ★ Bedford Av Kiosk"),
	)
	b.WriteString(radarTelemetry + "\n")

	// Amber CRT perspective scan grid
	grid1 := styles.FgBorder.Render("    / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / / /")
	grid2 := styles.FgPrimary.Render("    ═════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════")
	b.WriteString(grid1 + "\n")
	b.WriteString(grid2 + "\n")

	return b.String()
}
