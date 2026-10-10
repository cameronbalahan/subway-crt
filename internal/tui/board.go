package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"subway-crt/internal/transit"
)

// RenderDepartureBoard renders the top 1/3 Bedford Ave (L) commuter board
func RenderDepartureBoard(width int, snap transit.TransitSnapshot, styles Styles) string {
	if width < 80 {
		width = 80
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
	colWidth := (width - 8) / 2
	if colWidth < 46 {
		colWidth = 46
	}

	leftCol := renderRetroDepartureCard(
		"▲ TRACK 1: MANHATTAN-BOUND // 8 AV TERMINUS",
		snap.ManhattanTrains,
		colWidth,
		styles.Palette.Secondary,
		styles.Palette.Primary,
		styles,
	)

	rightCol := renderRetroDepartureCard(
		"▼ TRACK 2: CANARSIE-BOUND // ROCKAWAY PKWY",
		snap.CanarsieTrains,
		colWidth,
		styles.Palette.Primary,
		styles.Palette.Secondary,
		styles,
	)

	boardContent := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, "  ", rightCol)

	return lipgloss.JoinVertical(lipgloss.Left,
		alertLine,
		boardContent,
	)
}

func renderRetroDepartureCard(title string, departures []transit.TrainDeparture, width int, headerColor, accentColor lipgloss.Color, styles Styles) string {
	var b strings.Builder

	headerBar := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
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

	b.WriteString(styles.FgMuted.Render(fmt.Sprintf("  %-3s  %-14s  %-24s  %s", "SEQ", "COUNTDOWN", "LIVE TRAIN LOCATION", "ARRIVAL")) + "\n")
	b.WriteString(styles.FgBorder.Render("  " + strings.Repeat("┄", width-8)) + "\n")

	count := 0
	for _, d := range departures {
		if count >= 4 {
			break
		}
		count++

		// Styled outrun countdown pill
		var badge string
		if d.SecondsAway <= 45 {
			badge = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(styles.Palette.Alert).
				Blink(true).
				Padding(0, 1).
				Render("⚡ BOARDING ⚡")
		} else if d.SecondsAway <= 90 {
			badge = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(styles.Palette.Accent).
				Blink(true).
				Padding(0, 1).
				Render(fmt.Sprintf("▶ %2d SEC AWAY", d.SecondsAway))
		} else {
			badge = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(headerColor).
				Padding(0, 1).
				Render(fmt.Sprintf("  %2d MIN  ", d.MinutesAway))
		}

		statusDesc := d.CurrentStatus
		if len(statusDesc) > 22 {
			statusDesc = statusDesc[:22]
		}

		timeStr := d.ArrivalTime.Format("15:04:05")

		// Sequential indicator
		seqGlyph := fmt.Sprintf("#0%d", count)
		seqStyled := lipgloss.NewStyle().Bold(true).Foreground(accentColor).Render(seqGlyph)

		locStyled := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Render(fmt.Sprintf("%-22s", statusDesc))
		arrStyled := styles.FgMuted.Render(timeStr)

		row := fmt.Sprintf("  %s  %-14s  %s  %s",
			seqStyled,
			badge,
			locStyled,
			arrStyled,
		)
		b.WriteString(row + "\n")
	}

	card := styles.CardBox.Width(width).Render(b.String())
	return card
}
