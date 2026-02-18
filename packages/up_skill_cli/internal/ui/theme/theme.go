package theme

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Header             lipgloss.Style
	Prompt             lipgloss.Style
	Error              lipgloss.Style
	ListTitleBar       lipgloss.Style
	ListItem           lipgloss.Style
	ListTitle          lipgloss.Style
	ListItemHeader     lipgloss.Style
	ListSelectedHeader lipgloss.Style
	ListItemDesc       lipgloss.Style
	ListArrow          lipgloss.Style
}

func DefaultStyles() *Styles {
	primary := lipgloss.Color("#89dceb")
	accent := lipgloss.Color("#a6e3a1")
	muted := lipgloss.Color("#bac2de")
	warning := lipgloss.Color("#f9e2af")

	return &Styles{
		Header:             lipgloss.NewStyle().Foreground(primary).Bold(true),
		Prompt:             lipgloss.NewStyle().Foreground(accent).Bold(true),
		ListTitleBar:       lipgloss.NewStyle().Foreground(primary).Bold(true).Padding(0, 0, 1, 0),
		ListItem:           lipgloss.NewStyle().Padding(0, 0, 1, 0),
		ListArrow:          lipgloss.NewStyle().Foreground(accent).Bold(true),
		ListTitle:          lipgloss.NewStyle().Foreground(primary).Bold(true),
		ListItemHeader:     lipgloss.NewStyle().Bold(true).Padding(0, 0, 0, 2),
		ListSelectedHeader: lipgloss.NewStyle().Foreground(accent).Bold(true),
		ListItemDesc:       lipgloss.NewStyle().Foreground(muted).Italic(true).Padding(0, 0, 0, 2),
		Error:              lipgloss.NewStyle().Foreground(warning).Bold(true),
	}
}
