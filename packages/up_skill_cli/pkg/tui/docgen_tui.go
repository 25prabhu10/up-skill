package tui

import (
	textinput "github.com/25prabhu10/up-skill/internal/ui/text_input"

	tea "github.com/charmbracelet/bubbletea"
)

type docgenConfig struct {
	OutputPath  string
	Format      string
	FrontMatter bool
}

type dogenModel struct {
	config    docgenConfig
	pathInput *textinput.Model
}

func newDocgenModel(config *docgenConfig) dogenModel {
	pathInput := &textinput.Output{}
	input := textinput.InitialTextInputModelWithPlaceholder(pathInput, "Where to put the docs?", config.OutputPath, config.OutputPath)

	return dogenModel{
		config:    *config,
		pathInput: input,
	}
}

func (m dogenModel) Init() tea.Cmd {
	return nil
}

func (m dogenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m dogenModel) View() string {
	return "Docgen TUI"
}

func RunDocGenTUI(config *docgenConfig) error {
	p := tea.NewProgram(newDocgenModel(config))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
