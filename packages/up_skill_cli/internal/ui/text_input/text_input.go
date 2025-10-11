package textinput

import (
	"25prabhu10/up-skill/cli/internal/program"
	"25prabhu10/up-skill/cli/internal/ui/theme"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const defaultPlaceholder = "Inverse Binary Tree"

type (
	errMsg error
)

type Output struct {
	Output string
}

func (o *Output) update(val string) {
	o.Output = val
}

type model struct {
	textInput textinput.Model
	output    *Output
	err       error
	header    string
	exit      *bool
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if utf8.RuneCountInString(m.textInput.Value()) > 0 {
				m.output.update(m.textInput.Value())
			}
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			*m.exit = true
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		*m.exit = true
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("\n%s\n\n%s\n\n", m.header, m.textInput.View())
}

func (m model) Err() string {
	if m.err == nil {
		return ""
	}
	return m.err.Error()
}

func InitialTextInputModel(output *Output, header string, initialStr string) *model {
	ti := newTextInputModel(defaultPlaceholder, initialStr)
	exit := false

	return &model{
		textInput: ti,
		output:    output,
		header:    theme.Default.Title.Render(header),
		err:       nil,
		exit:      &exit,
	}
}

func InitialTextInputModelWithPlaceholder(output *Output, header string, initialStr string, placeholder string) *model {
	ti := newTextInputModel(placeholder, initialStr)
	exit := false

	return &model{
		textInput: ti,
		output:    output,
		header:    theme.Default.Title.Render(header),
		err:       nil,
		exit:      &exit,
	}
}

func InitialTextInputTopicModel(output *Output, header string, program *program.Topic, initialStr string) *model {
	ti := newTextInputModel(defaultPlaceholder, initialStr)

	return &model{
		textInput: ti,
		output:    output,
		header:    theme.Default.Title.Render(header),
		err:       nil,
		exit:      &program.Exit,
	}
}

func InitialTextInputTopicModelWithPlaceholder(output *Output, header string, program *program.Topic, initialStr string, placeholder string) *model {
	ti := newTextInputModel(placeholder, initialStr)

	return &model{
		textInput: ti,
		output:    output,
		header:    theme.Default.Title.Render(header),
		err:       nil,
		exit:      &program.Exit,
	}
}

func CreateErrorInputModel(err error) *model {
	ti := newTextInputModel("", "")
	exit := true

	return &model{
		textInput: ti,
		output:    nil,
		header:    "",
		err:       errors.New(theme.Default.Error.Render(err.Error())),
		exit:      &exit,
	}
}

func newTextInputModel(placeholder string, initialStr string) textinput.Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	ti.PromptStyle = theme.Default.InputPrompt
	ti.Placeholder = placeholder

	if initialStr != "" {
		ti.SetValue(initialStr)
	}

	return ti
}
