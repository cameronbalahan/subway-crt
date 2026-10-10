package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"subway-crt/internal/transit"
)

type TickMsg time.Time

type RootModel struct {
	width        int
	height       int
	snapshotChan <-chan transit.TransitSnapshot
	snapshot     transit.TransitSnapshot
	lastTick     time.Time
	styles       Styles
}

func NewRootModel(snapshotChan <-chan transit.TransitSnapshot) RootModel {
	return RootModel{
		snapshotChan: snapshotChan,
		width:        240, // 1080p full screen default width
		height:       67,  // 1080p full screen default height
		lastTick:     time.Now(),
		styles:       MakeStyles(GetDefaultPalette()),
	}
}

func (m RootModel) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case TickMsg:
		m.lastTick = time.Time(msg)

		// Drain latest snapshot if available on channel
		select {
		case snap := <-m.snapshotChan:
			m.snapshot = snap
		default:
			// Decrement local seconds countdown between network ticks
			for i := range m.snapshot.ManhattanTrains {
				if m.snapshot.ManhattanTrains[i].SecondsAway > 0 {
					m.snapshot.ManhattanTrains[i].SecondsAway--
					m.snapshot.ManhattanTrains[i].MinutesAway = m.snapshot.ManhattanTrains[i].SecondsAway / 60
				}
			}
			for i := range m.snapshot.CanarsieTrains {
				if m.snapshot.CanarsieTrains[i].SecondsAway > 0 {
					m.snapshot.CanarsieTrains[i].SecondsAway--
					m.snapshot.CanarsieTrains[i].MinutesAway = m.snapshot.CanarsieTrains[i].SecondsAway / 60
				}
			}
		}

		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return TickMsg(t)
		})
	}

	return m, nil
}

func (m RootModel) View() string {
	w := m.width
	if w < 160 {
		w = 160
	}

	// 1. Amber Monochrome CRT Command Center Header
	nowStr := m.lastTick.Format("15:04:05")
	feedMode := m.snapshot.FeedMode
	if feedMode == "" {
		feedMode = "LIVE_GTFS_RT"
	}

	kioskTitle := " ◆ SUBWAY-CRT // BEDFORD AV (L) 1080p RADAR // AMBER MONOCHROME P3 "
	statusInfo := fmt.Sprintf(" [FEED: %s] [%s EST] ", feedMode, nowStr)

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		m.styles.KioskHeader.Render(kioskTitle),
		" ",
		m.styles.Clock.Render(statusInfo),
	)

	// 2. Top Half: Widescreen Departure Board with fixed-width columns
	topView := RenderDepartureBoard(w, m.snapshot, m.styles)

	// 3. Amber Divider
	divider := lipgloss.NewStyle().
		Foreground(m.styles.Palette.Border).
		Render("───◆ " + strings.Repeat("━", w-12) + " ◆───")

	// 4. Bottom Half: Centered, widened Schematic Map
	bottomView := RenderSchematicMap(w, m.height-18, m.snapshot, m.styles)

	// 5. Clean Footer Bar (No palette changing, no engine credits)
	badgeQuit := m.styles.KeyBadge.Render("q")
	badgeKiosk := m.styles.KeyBadge.Render("1080p")

	footer := fmt.Sprintf(
		"  %s Quit  ·  %s Console  ·  BEDFORD AV (L) PASSENGER RADAR KIOSK",
		badgeQuit,
		badgeKiosk,
	)
	footerStyled := m.styles.FooterBar.Render(footer)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"\n",
		topView,
		"\n",
		divider,
		"\n",
		bottomView,
		"\n",
		footerStyled,
	)
}

func (m RootModel) ViewWithSnapshot(snap transit.TransitSnapshot) string {
	m.snapshot = snap
	return m.View()
}
