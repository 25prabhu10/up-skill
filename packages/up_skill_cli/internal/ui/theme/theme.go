package theme

import "github.com/charmbracelet/lipgloss"

const (
	colorSurface    = "#1e1e2e"
	colorText       = "#cdd6f4"
	colorSubtleText = "#a6adc8"
	colorAccent     = "#89b4fa"
	colorError      = "#f38ba8"
)

type Styles struct {
	Title              lipgloss.Style
	ListTitle          lipgloss.Style
	Error              lipgloss.Style
	Help               lipgloss.Style
	InputPrompt        lipgloss.Style
	InputText          lipgloss.Style
	ListItem           lipgloss.Style
	ListHeader         lipgloss.Style
	ListItemDesc       lipgloss.Style
	ListSelectedItem   lipgloss.Style
	ListSelectedHeader lipgloss.Style
	ListArrow          lipgloss.Style
	ListSelectedDesc   lipgloss.Style
}

var Default = Styles{
	Title:              lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Bold(true).Padding(0, 0, 0, 2),
	ListTitle:          lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Bold(true),
	InputText:          lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)),
	InputPrompt:        lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent)).Bold(true),
	ListItem:           lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Padding(0, 0, 1, 2),
	ListHeader:         lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Bold(true),
	ListItemDesc:       lipgloss.NewStyle().Foreground(lipgloss.Color(colorSubtleText)).Italic(true),
	ListArrow:          lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent)).Bold(true),
	ListSelectedItem:   lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent)).Padding(0, 0, 1, 0),
	ListSelectedHeader: lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent)).Bold(true),
	ListSelectedDesc:   lipgloss.NewStyle().Foreground(lipgloss.Color(colorSubtleText)).Padding(0, 0, 0, 2).Italic(true),
	Help:               lipgloss.NewStyle().Foreground(lipgloss.Color(colorSubtleText)).PaddingLeft(2),
	Error:              lipgloss.NewStyle().Foreground(lipgloss.Color(colorError)).Bold(true),
}
