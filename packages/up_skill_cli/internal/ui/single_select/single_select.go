package single_select

import (
	"25prabhu10/up-skill/cli/internal/options"
	"25prabhu10/up-skill/cli/internal/ui/theme"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

const (
	listWidth     = 24
	minListHeight = 5
)

type Output struct {
	Output string
}

func (o *Output) update(val string) {
	o.Output = val
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		title, desc string
	)

	if i, ok := listItem.(options.Item); ok {
		title = i.Title()
		desc = i.Description()
	} else {
		return
	}

	if m.Width() <= 0 {
		// short-circuit
		return
	}

	textWidth := m.Width() - theme.Default.ListItem.GetPaddingLeft() - theme.Default.ListItem.GetPaddingRight()
	var str string
	if index == m.Index() {
		title = theme.Default.ListSelectedHeader.Render(ansi.Truncate(title, textWidth, "…"))
		desc = theme.Default.ListSelectedDesc.Render(desc)
		arrow := theme.Default.ListArrow.Render(">")
		str = theme.Default.ListSelectedItem.Render(fmt.Sprintf("%s %d. %s\n   %s", arrow, index+1, title, desc))
	} else {
		title = theme.Default.ListHeader.Render(ansi.Truncate(title, textWidth, "…"))
		desc = theme.Default.ListItemDesc.Render(desc)
		str = theme.Default.ListItem.Render(fmt.Sprintf("%d. %s\n   %s", index+1, title, desc))
	}

	fmt.Fprint(w, str)
}

type model struct {
	list   list.Model
	output *Output
	err    error
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(options.Item)
			if ok {
				m.output.update(i.Title())
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	return m.list.View()
}

func InitializeModel(userOptions *[]list.Item, title *string, output *Output) *model {
	height := len(*userOptions) * 5

	m := model{
		list:   list.New(*userOptions, itemDelegate{}, listWidth, height),
		output: output,
	}

	m.list.Title = *title
	m.list.SetShowStatusBar(false)

	m.list.Styles.Title = theme.Default.ListTitle

	return &m
}
