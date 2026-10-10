package tui

import (
	"fmt"
	"runtime"
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
	currentTheme ThemeType
	styles       Styles
}

func NewRootModel(snapshotChan <-chan transit.TransitSnapshot) RootModel {
	theme := ThemeOutrun
	return RootModel{
		snapshotChan: snapshotChan,
		width:        240,
		height:       67,
		lastTick:     time.Now(),
		currentTheme: theme,
		styles:       MakeStyles(GetPalette(theme)),
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
		case "p", "tab":
			// Cycle themes: Outrun -> Phosphor -> Amber -> Outrun
			switch m.currentTheme {
			case ThemeOutrun:
				m.currentTheme = ThemePhosphor
			case ThemePhosphor:
				m.currentTheme = ThemeAmber
			default:
				m.currentTheme = ThemeOutrun
			}
			m.styles = MakeStyles(GetPalette(m.currentTheme))
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

func (m RootModel) SetTheme(theme ThemeType) RootModel {
	m.currentTheme = theme
	m.styles = MakeStyles(GetPalette(theme))
	return m
}

func (m RootModel) View() string {
	w := m.width
	if w < 100 {
		w = 100
	}

	// 1. Synthwave Command Center Header
	nowStr := m.lastTick.Format("15:04:05")
	feedMode := m.snapshot.FeedMode
	if feedMode == "" {
		feedMode = "LIVE_GTFS_RT"
	}

	kioskTitle := fmt.Sprintf(" ◆ SUBWAY-CRT // BEDFORD AV (L) 1080p RADAR // %s ", m.styles.Palette.Name)
	statusInfo := fmt.Sprintf(" [FEED: %s] [%s EST] ", feedMode, nowStr)

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		m.styles.KioskHeader.Render(kioskTitle),
		" ",
		m.styles.Clock.Render(statusInfo),
	)

	// 2. Top Half: Outrun Departure Board
	topView := RenderDepartureBoard(w, m.snapshot, m.styles)

	// 3. Neon Divider with cyber glyph accents
	divider := lipgloss.NewStyle().
		Foreground(m.styles.Palette.Border).
		Render("───◆ " + strings.Repeat("━", w-12) + " ◆───")

	// 4. Bottom Half: Outrun Schematic Map
	bottomView := RenderSchematicMap(w, m.height-18, m.snapshot, m.styles)

	// 5. Interactive Outrun Keymap Footer
	badgeQuit := m.styles.KeyBadge.Render("q")
	badgeTheme := m.styles.KeyBadge.Render("p / tab")
	badgeKiosk := m.styles.KeyBadge.Render("1080p")
	goVer := strings.TrimPrefix(runtime.Version(), "go")
	badgePureGo := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#00ADD8")).Padding(0, 1).Render("Pure Go " + goVer)

	footer := fmt.Sprintf(
		"  %s Quit  ·  %s Palette (%s)  ·  %s Console  ·  %s  ·  Official MTA GTFS-RT Protobuf Engine · Charm Lip Gloss ",
		badgeQuit,
		badgeTheme,
		m.styles.Palette.Name,
		badgeKiosk,
		badgePureGo,
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
