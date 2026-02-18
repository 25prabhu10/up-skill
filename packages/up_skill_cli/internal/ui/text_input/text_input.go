package textinput

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/25prabhu10/up-skill/internal/program"
	"github.com/25prabhu10/up-skill/internal/ui/theme"

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

type Model struct {
	textInput textinput.Model
	output    *Output
	err       error
	header    string
	exit      *bool
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m Model) View() string {
	return fmt.Sprintf("\n%s\n\n%s\n\n", m.header, m.textInput.View())
}

func (m Model) Err() string {
	if m.err == nil {
		return ""
	}
	return m.err.Error()
}

func InitialTextInputModel(output *Output, header string, initialStr string) *Model {
	ti := newTextInputModel(defaultPlaceholder, initialStr)
	exit := false

	return &Model{
		textInput: ti,
		output:    output,
		header:    theme.DefaultStyles().Header.Render(header),
		err:       nil,
		exit:      &exit,
	}
}

func InitialTextInputModelWithPlaceholder(output *Output, header string, initialStr string, placeholder string) *Model {
	ti := newTextInputModel(placeholder, initialStr)
	exit := false

	return &Model{
		textInput: ti,
		output:    output,
		header:    theme.DefaultStyles().Header.Render(header),
		err:       nil,
		exit:      &exit,
	}
}

func InitialTextInputTopicModel(output *Output, header string, program *program.Program, initialStr string) *Model {
	ti := newTextInputModel(defaultPlaceholder, initialStr)

	return &Model{
		textInput: ti,
		output:    output,
		header:    theme.DefaultStyles().Header.Render(header),
		err:       nil,
		exit:      nil,
	}
}

func InitialTextInputTopicModelWithPlaceholder(output *Output, header string, program *program.Program, initialStr string, placeholder string) *Model {
	ti := newTextInputModel(placeholder, initialStr)

	return &Model{
		textInput: ti,
		output:    output,
		header:    theme.DefaultStyles().Header.Render(header),
		err:       nil,
		exit:      nil,
	}
}

func CreateErrorInputModel(err error) *Model {
	ti := newTextInputModel("", "")
	exit := true

	return &Model{
		textInput: ti,
		output:    nil,
		header:    "",
		err:       errors.New(theme.DefaultStyles().Error.Render(err.Error())),
		exit:      &exit,
	}
}

func newTextInputModel(placeholder string, initialStr string) textinput.Model {
	ti := textinput.New()
	ti.CharLimit = 512
	ti.Placeholder = placeholder

	if initialStr != "" {
		ti.SetValue(initialStr)
	}

	styles := theme.DefaultStyles()
	ti.PromptStyle = styles.Prompt
	ti.Focus()
	return ti
}
