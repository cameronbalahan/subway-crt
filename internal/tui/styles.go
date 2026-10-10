package tui

import "github.com/charmbracelet/lipgloss"

type ThemeType string

const (
	ThemeOutrun   ThemeType = "outrun"
	ThemePhosphor ThemeType = "phosphor"
	ThemeAmber    ThemeType = "amber"
)

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
	// Outrun Synthwave Palette (Neon Magenta / Cyber Cyan / Electric Gold)
	PaletteOutrun = ThemePalette{
		Name:       "SYNTHWAVE OUTRUN",
		Primary:    lipgloss.Color("#FF007F"), // Neon Pink / Magenta
		Secondary:  lipgloss.Color("#00E5FF"), // Cyber Cyan
		Accent:     lipgloss.Color("#FFE600"), // Electric Yellow
		Alert:      lipgloss.Color("#FF0055"), // Critical Red
		Border:     lipgloss.Color("#9D00FF"), // Neon Violet
		Background: lipgloss.Color("#07070F"), // Deep Void
		Surface:    lipgloss.Color("#120F24"), // Cyber Surface
		Text:       lipgloss.Color("#FFFFFF"),
		Muted:      lipgloss.Color("#8888AA"),
		Track:      lipgloss.Color("#00F0FF"),
		River:      lipgloss.Color("#00B4D8"),
	}

	// CRT Phosphor Green Palette (Classic Matrix / Radar)
	PalettePhosphor = ThemePalette{
		Name:       "PHOSPHOR CRT (P1)",
		Primary:    lipgloss.Color("#00FF66"), // CRT Green
		Secondary:  lipgloss.Color("#70FFB8"), // Mint
		Accent:     lipgloss.Color("#00FFCC"),
		Alert:      lipgloss.Color("#FF3366"),
		Border:     lipgloss.Color("#008833"),
		Background: lipgloss.Color("#0B0B10"),
		Surface:    lipgloss.Color("#0E1612"),
		Text:       lipgloss.Color("#00FF66"),
		Muted:      lipgloss.Color("#009944"),
		Track:      lipgloss.Color("#00FF66"),
		River:      lipgloss.Color("#00CC88"),
	}

	// CRT Amber Monochrome (Classic amber terminal)
	PaletteAmber = ThemePalette{
		Name:       "AMBER MONOCHROME (P3)",
		Primary:    lipgloss.Color("#FFB000"), // Warm Amber
		Secondary:  lipgloss.Color("#FFD066"),
		Accent:     lipgloss.Color("#FF9900"),
		Alert:      lipgloss.Color("#FF2200"),
		Border:     lipgloss.Color("#996600"),
		Background: lipgloss.Color("#0E0B05"),
		Surface:    lipgloss.Color("#1A1408"),
		Text:       lipgloss.Color("#FFB000"),
		Muted:      lipgloss.Color("#997722"),
		Track:      lipgloss.Color("#FFB000"),
		River:      lipgloss.Color("#CC8800"),
	}
)

func GetPalette(t ThemeType) ThemePalette {
	switch t {
	case ThemePhosphor:
		return PalettePhosphor
	case ThemeAmber:
		return PaletteAmber
	default:
		return PaletteOutrun
	}
}

// Outrun Styles Generator using current theme palette
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

	BadgeACE  lipgloss.Style
	Badge123  lipgloss.Style
	BadgeBDFM lipgloss.Style
	BadgeNQRW lipgloss.Style
	Badge456  lipgloss.Style
	BadgeG    lipgloss.Style
	BadgeJMZ  lipgloss.Style
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

		BadgeACE:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#0039A6")).Padding(0, 1),
		Badge123:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#EE352E")).Padding(0, 1),
		BadgeBDFM: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#FF6319")).Padding(0, 1),
		BadgeNQRW: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FCCC0A")).Padding(0, 1),
		Badge456:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#00933C")).Padding(0, 1),
		BadgeG:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#6CBE45")).Padding(0, 1),
		BadgeJMZ:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#996633")).Padding(0, 1),

		KioskHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Text).
			Background(p.Primary).
			Padding(0, 2),

		KioskTag: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Background).
			Background(p.Secondary).
			Padding(0, 1),

		Clock: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Secondary).
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
			Foreground(p.Primary),

		BeaconS: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Secondary),

		FooterBar: lipgloss.NewStyle().
			Foreground(p.Muted).
			Background(p.Background),

		KeyBadge: lipgloss.NewStyle().
			Bold(true).
			Foreground(p.Background).
			Background(p.Secondary).
			Padding(0, 1),
	}
}
