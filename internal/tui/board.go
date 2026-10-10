package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"subway-crt/internal/transit"
)

// RenderDepartureBoard renders the top 1/3 Bedford Ave (L) commuter board
// Optimized for wide 1080p full-screen display with rock-solid column alignment
func RenderDepartureBoard(width int, snap transit.TransitSnapshot, styles Styles) string {
	if width < 120 {
		width = 120
	}

	// 1. Critical Alert Banner or Status Line
	var alertLine string
	if snap.IsStale {
		elapsed := int(time.Since(snap.LastSuccessTime).Seconds())
		alertLine = styles.AlertCritical.Width(width - 4).Render(
			fmt.Sprintf("⚡ ░▒▓ WARNING: LIVE FEED DEGRADED (%ds ago) · CBTC KINEMATIC SIMULATION ACTIVE ▓▒░ ⚡", elapsed),
		)
	} else {
		hasCrit := false
		for _, a := range snap.ActiveAlerts {
			if a.Severity == "CRITICAL" && a.AffectsL {
				alertLine = styles.AlertCritical.Width(width - 4).Render(
					fmt.Sprintf("▲ ░▒▓ CRITICAL MTA ADVISORY: %s ▓▒░ ▲", a.Summary),
				)
				hasCrit = true
				break
			}
		}
		if !hasCrit && len(snap.ActiveAlerts) > 0 {
			a := snap.ActiveAlerts[0]
			if a.Severity == "WARNING" {
				alertLine = styles.AlertNotice.Width(width - 4).Render(
					fmt.Sprintf("◆ NOTICE: %s · OPERATING ALTERNATE PATTERN", a.Summary),
				)
			} else {
				alertLine = styles.NominalBar.Width(width - 4).Render(
					fmt.Sprintf("⚡ STATUS: %s · DUAL-TRACK AUTOMATIC TRAIN CONTROL ACTIVE", a.Summary),
				)
			}
		}
	}

	// 2. Split Track Departure Cards: Manhattan (Left) vs Canarsie (Right)
	// Full wide screen spacing: split available width evenly with padding
	cardWidth := (width - 8) / 2
	if cardWidth < 68 {
		cardWidth = 68
	}

	leftCol := renderFixedDepartureCard(
		"▲ TRACK 1: MANHATTAN-BOUND // 8 AV TERMINUS",
		snap.ManhattanTrains,
		cardWidth,
		styles.Palette.Primary,
		styles.Palette.Accent,
		styles,
	)

	rightCol := renderFixedDepartureCard(
		"▼ TRACK 2: CANARSIE-BOUND // ROCKAWAY PKWY",
		snap.CanarsieTrains,
		cardWidth,
		styles.Palette.Secondary,
		styles.Palette.Primary,
		styles,
	)

	boardContent := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, "  ", rightCol)

	if alertLine != "" {
		return lipgloss.JoinVertical(lipgloss.Left,
			alertLine,
			boardContent,
		)
	}
	return boardContent
}

// renderFixedDepartureCard guarantees strict tabular monospace column widths
// Columns:
//   SEQ (4 chars: " #01 ")
//   COUNTDOWN (14 chars: "[  3 MIN     ]" / "[ 42 SEC AWAY]" / "[ BOARDING   ]")
//   DESTINATION (22 chars)
//   LIVE TRAIN LOCATION (26 chars)
//   ARRIVAL TIME (12 chars: "16:54:15 EDT")
func renderFixedDepartureCard(title string, departures []transit.TrainDeparture, width int, headerColor, accentColor lipgloss.Color, styles Styles) string {
	var b strings.Builder

	headerBar := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Palette.Background).
		Background(headerColor).
		Padding(0, 2).
		Width(width - 4).
		Render(title)

	b.WriteString(headerBar)
	b.WriteString("\n")

	if len(departures) == 0 {
		b.WriteString("\n")
		b.WriteString(styles.StatusText.Render("  [!] NO SCHEDULED TRAIN RUNS DETECTED IN CORRIDOR\n"))
		b.WriteString("\n")
		card := styles.CardBox.Width(width).Render(b.String())
		return card
	}

	// Strictly fixed column header
	headerRow := fmt.Sprintf("  %-5s %-16s %-22s %-26s %s",
		"SEQ", "COUNTDOWN", "DESTINATION", "LIVE TRAIN LOCATION", "ARRIVAL")
	b.WriteString(styles.FgMuted.Render(headerRow) + "\n")
	b.WriteString(styles.FgBorder.Render("  "+strings.Repeat("┄", width-8)) + "\n")

	count := 0
	for _, d := range departures {
		if count >= 4 {
			break
		}
		count++

		// 1. Strictly fixed-width countdown box (16 chars display width)
		var countdownPill string
		if d.SecondsAway <= 45 {
			countdownPill = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(styles.Palette.Alert).
				Blink(true).
				Render(" [ BOARDING ]   ")
		} else if d.SecondsAway <= 90 {
			countdownPill = lipgloss.NewStyle().
				Bold(true).
				Foreground(styles.Palette.Background).
				Background(styles.Palette.Accent).
				Blink(true).
				Render(fmt.Sprintf(" [ %2d SEC AWAY ]", d.SecondsAway))
		} else {
			countdownPill = lipgloss.NewStyle().
				Bold(true).
				Foreground(styles.Palette.Background).
				Background(styles.Palette.Primary).
				Render(fmt.Sprintf(" [   %2d MIN    ] ", d.MinutesAway))
		}

		// 2. Fixed Destination (truncated or padded to exactly 22 chars)
		dest := d.Destination
		if len(dest) > 22 {
			dest = dest[:22]
		}
		destFormatted := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.Palette.Text).
			Render(fmt.Sprintf("%-22s", dest))

		// 3. Fixed Live Location (truncated or padded to exactly 26 chars)
		loc := d.CurrentStatus
		if len(loc) > 26 {
			loc = loc[:26]
		}
		locFormatted := lipgloss.NewStyle().
			Foreground(styles.Palette.Secondary).
			Render(fmt.Sprintf("%-26s", loc))

		// 4. Fixed Arrival Time (12 chars: "16:54:15 EDT")
		timeStr := d.ArrivalTime.Format("15:04:05 EDT")
		arrFormatted := styles.FgMuted.Render(fmt.Sprintf("%-12s", timeStr))

		// 5. Sequence indicator (5 chars)
		seqFormatted := lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Render(fmt.Sprintf("#0%d  ", count))

		// Assemble row with strict spacing: fields will NEVER shift horizontally
		row := fmt.Sprintf("  %s%s %s %s %s",
			seqFormatted,
			countdownPill,
			destFormatted,
			locFormatted,
			arrFormatted,
		)
		b.WriteString(row + "\n")
	}

	card := styles.CardBox.Width(width).Render(b.String())
	return card
}
