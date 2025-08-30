package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme colors
var (
	PrimaryColor   = lipgloss.Color("#7C3AED")
	SecondaryColor = lipgloss.Color("#059669")
	AccentColor    = lipgloss.Color("#DC2626")
	SuccessColor   = lipgloss.Color("#16A34A")
	WarningColor   = lipgloss.Color("#CA8A04")
	ErrorColor     = lipgloss.Color("#DC2626")
	InfoColor      = lipgloss.Color("#0284C7")

	TextColor    = lipgloss.Color("#F3F4F6")
	MutedColor   = lipgloss.Color("#9CA3AF")
	BorderColor  = lipgloss.Color("#374151")
	BgColor      = lipgloss.Color("#111827")
	SurfaceColor = lipgloss.Color("#1F2937")
)

// Theme system with consistent styling
type Theme struct {
	// Base styles
	Base    lipgloss.Style
	Surface lipgloss.Style
	Border  lipgloss.Style

	// Text styles
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Body     lipgloss.Style
	Muted    lipgloss.Style

	// Interactive styles
	Button       lipgloss.Style
	ButtonActive lipgloss.Style
	Input        lipgloss.Style
	InputFocused lipgloss.Style

	// Status styles
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style
	Info    lipgloss.Style

	// Layout styles
	Panel          lipgloss.Style
	Card           lipgloss.Style
	List           lipgloss.Style
	ListItem       lipgloss.Style
	ListItemActive lipgloss.Style

	// Navigation styles
	TabActive   lipgloss.Style
	TabInactive lipgloss.Style
	Breadcrumb  lipgloss.Style

	// Progress and status
	ProgressBar   lipgloss.Style
	StatusGood    lipgloss.Style
	StatusBad     lipgloss.Style
	StatusUnknown lipgloss.Style
}

func NewTheme() *Theme {
	return &Theme{
		// Base styles
		Base: lipgloss.NewStyle().
			Foreground(TextColor).
			Background(BgColor),

		Surface: lipgloss.NewStyle().
			Background(SurfaceColor).
			Padding(1, 2),

		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor),

		// Text styles
		Title: lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			MarginBottom(1),

		Subtitle: lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true),

		Body: lipgloss.NewStyle().
			Foreground(TextColor),

		Muted: lipgloss.NewStyle().
			Foreground(MutedColor),

		// Interactive styles
		Button: lipgloss.NewStyle().
			Foreground(TextColor).
			Background(SurfaceColor).
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor),

		ButtonActive: lipgloss.NewStyle().
			Foreground(BgColor).
			Background(PrimaryColor).
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor),

		Input: lipgloss.NewStyle().
			Foreground(TextColor).
			Background(SurfaceColor).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor),

		InputFocused: lipgloss.NewStyle().
			Foreground(TextColor).
			Background(SurfaceColor).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor),

		// Status styles
		Success: lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true),

		Warning: lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true),

		Info: lipgloss.NewStyle().
			Foreground(InfoColor).
			Bold(true),

		// Layout styles
		Panel: lipgloss.NewStyle().
			Background(SurfaceColor).
			Padding(1, 2).
			Margin(1, 0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor),

		Card: lipgloss.NewStyle().
			Background(SurfaceColor).
			Padding(1, 2).
			Margin(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor),

		List: lipgloss.NewStyle().
			Background(SurfaceColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor),

		ListItem: lipgloss.NewStyle().
			Foreground(TextColor).
			Padding(0, 1),

		ListItemActive: lipgloss.NewStyle().
			Foreground(BgColor).
			Background(PrimaryColor).
			Padding(0, 1).
			Bold(true),

		// Navigation styles
		TabActive: lipgloss.NewStyle().
			Foreground(BgColor).
			Background(PrimaryColor).
			Padding(0, 2).
			Bold(true),

		TabInactive: lipgloss.NewStyle().
			Foreground(MutedColor).
			Background(SurfaceColor).
			Padding(0, 2),

		Breadcrumb: lipgloss.NewStyle().
			Foreground(MutedColor).
			MarginBottom(1),

		// Progress and status
		ProgressBar: lipgloss.NewStyle().
			Background(SurfaceColor).
			Foreground(PrimaryColor),

		StatusGood: lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true),

		StatusBad: lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true),

		StatusUnknown: lipgloss.NewStyle().
			Foreground(MutedColor).
			Bold(true),
	}
}

// Responsive layout utilities
func AdaptiveWidth(width int, minWidth int) int {
	if width < minWidth {
		return minWidth
	}
	return width
}

func AdaptiveHeight(height int, minHeight int) int {
	if height < minHeight {
		return minHeight
	}
	return height
}

// Progress bar rendering
func RenderProgressBar(progress float64, width int, theme *Theme) string {
	if width <= 0 {
		return ""
	}

	filled := int(progress * float64(width))
	if filled > width {
		filled = width
	}

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return theme.ProgressBar.Render(bar)
}

// Status indicator rendering
func RenderStatusIndicator(status string, theme *Theme) string {
	switch status {
	case "active", "running", "healthy", "connected", "online":
		return theme.StatusGood.Render("● " + status)
	case "inactive", "stopped", "unhealthy", "disconnected", "offline", "error", "failed":
		return theme.StatusBad.Render("● " + status)
	default:
		return theme.StatusUnknown.Render("● " + status)
	}
}

// Box rendering with adaptive sizing
func RenderBox(content string, title string, width, height int, theme *Theme) string {
	style := theme.Panel.
		Width(AdaptiveWidth(width-4, 20)).
		Height(AdaptiveHeight(height-2, 5))

	if title != "" {
		style = style.BorderTop(true).BorderTopForeground(PrimaryColor)
		content = theme.Subtitle.Render(title) + "\n" + content
	}

	return style.Render(content)
}
