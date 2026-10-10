package tui

import "github.com/charmbracelet/lipgloss"

type ThemePalette struct {
	Name        string
	Primary     lipgloss.Color
	Secondary   lipgloss.Color
	Accent      lipgloss.Color
	Alert       lipgloss.Color
	Border      lipgloss.Color
	Background  lipgloss.Color
	Surface     lipgloss.Color
	Text        lipgloss.Color
	Muted       lipgloss.Color
	Track       lipgloss.Color
	River       lipgloss.Color
}

var (
	// CRT Amber Monochrome (P3 Phosphor terminal)
	PaletteAmber = ThemePalette{
		Name:       "AMBER MONOCHROME (P3)",
		Primary:    lipgloss.Color("#FFB000"), // Warm Amber Phosphor
		Secondary:  lipgloss.Color("#FFC83B"), // Bright Amber Highlight
		Accent:     lipgloss.Color("#FFE066"), // Intense Amber Pip
		Alert:      lipgloss.Color("#FF3300"), // Red-Amber Critical Alert
		Border:     lipgloss.Color("#805300"), // Deep Amber Wireframe
		Background: lipgloss.Color("#0C0802"), // Black Phosphor CRT Tube
		Surface:    lipgloss.Color("#1A1104"), // Inner Bezel Shadow
		Text:       lipgloss.Color("#FFB000"), // Amber Monospace Text
		Muted:      lipgloss.Color("#A67200"), // Dim Amber
		Track:      lipgloss.Color("#FFB000"), // Amber Track Line
		River:      lipgloss.Color("#805300"), // River Barrier
	}
)

func GetDefaultPalette() ThemePalette {
	return PaletteAmber
}

type Styles struct {
	Palette ThemePalette

	KioskHeader    lipgloss.Style
	KioskTag       lipgloss.Style
	Clock          lipgloss.Style
	CardBox        lipgloss.Style
	CardHeader     lipgloss.Style
	CountdownHuge  lipgloss.Style
	CountdownBlink lipgloss.Style
	StatusText     lipgloss.Style
	AlertCritical  lipgloss.Style
	AlertNotice    lipgloss.Style
	NominalBar     lipgloss.Style
	TrackLine      lipgloss.Style
	StationPip     lipgloss.Style
	BedfordPip     lipgloss.Style
	BeaconN        lipgloss.Style
	BeaconS        lipgloss.Style
	FooterBar      lipgloss.Style
	KeyBadge       lipgloss.Style

	FgPrimary   lipgloss.Style
	FgSecondary lipgloss.Style
	FgAccent    lipgloss.Style
	FgMuted     lipgloss.Style
	FgRiver     lipgloss.Style
	FgBorder    lipgloss.Style
	FgText      lipgloss.Style

	BadgeAmberMain lipgloss.Style
	BadgeAmberDim  lipgloss.Style
}

func MakeStyles(p ThemePalette) Styles {
	return Styles{
		Palette: p,

		FgPrimary:   lipgloss.NewStyle().Foreground(p.Primary),
		FgSecondary: lipgloss.NewStyle().Foreground(p.Secondary),
		FgAccent:    lipgloss.NewStyle().Foreground(p.Accent),
		FgMuted:     lipgloss.NewStyle().Foreground(p.Muted),
		FgRiver:     lipgloss.NewStyle().Foreground(p.River),
		FgBorder:    lipgloss.NewStyle().Foreground(p.Border),
		FgText:      lipgloss.NewStyle().Foreground(p.Text),

		BadgeAmberMain: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Background).
			Background(p.Primary).
			Padding(0, 1),

		BadgeAmberDim: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Primary).
			Background(p.Surface).
			Padding(0, 1),

		KioskHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Background).
			Background(p.Primary).
			Padding(0, 2),

		KioskTag: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Background).
			Background(p.Secondary).
			Padding(0, 1),

		Clock: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Primary).
			Background(p.Surface).
			Padding(0, 2),

		CardBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(p.Border).
			Background(p.Surface).
			Padding(0, 1),

		CardHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Secondary),

		CountdownHuge: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Primary),

		CountdownBlink: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Accent).
			Blink(true),

		StatusText: lipgloss.NewStyle().
			Foreground(p.Muted),

		AlertCritical: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(p.Alert).
			Padding(0, 2),

		AlertNotice: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Background).
			Background(p.Accent).
			Padding(0, 2),

		NominalBar: lipgloss.NewStyle().
			Foreground(p.Secondary).
			Background(p.Surface).
			Padding(0, 2),

		TrackLine: lipgloss.NewStyle().
			Foreground(p.Track),

		StationPip: lipgloss.NewStyle().
			Foreground(p.Text),

		BedfordPip: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Accent),

		BeaconN: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Accent),

		BeaconS: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Secondary),

		FooterBar: lipgloss.NewStyle().
			Foreground(p.Muted).
			Background(p.Background),

		KeyBadge: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Background).
			Background(p.Primary).
			Padding(0, 1),
	}
}
