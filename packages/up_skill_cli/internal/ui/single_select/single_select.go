package single_select

import (
	"fmt"
	"io"
	"strings"

	"github.com/25prabhu10/up-skill/internal/ui/theme"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

const (
	listWidth = 24
)

type Option string

func (o Option) Title() string {
	return string(o)
}

func (o Option) Description() string {
	return ""
}

type Item struct {
	title Option
	desc  string
}

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
	var bodyLines strings.Builder

	var (
		title, desc string
	)

	// if i, ok := listItem.(Item); ok {
	// 	title = i.Title()
	// 	desc = i.Description()
	// } else {
	// 	return
	// }

	if m.Width() <= 0 {
		// short-circuit
		return
	}

	styles := theme.DefaultStyles()

	textWidth := m.Width() - styles.ListItem.GetPaddingLeft() - styles.ListItem.GetPaddingRight()
	if index == m.Index() {
		bodyLines.WriteString(styles.ListArrow.Render("> "))
		bodyLines.WriteString(styles.ListSelectedHeader.Render(ansi.Truncate(title, textWidth, "…")))
		bodyLines.WriteString("\n")
		bodyLines.WriteString(styles.ListItemDesc.Render(desc))
	} else {
		bodyLines.WriteString(styles.ListItemHeader.Render(ansi.Truncate(title, textWidth, "…")))
		bodyLines.WriteString("\n")
		bodyLines.WriteString(styles.ListItemDesc.Render(desc))
	}

	fmt.Fprint(w, styles.ListItem.Render(bodyLines.String()))
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
			// i, ok := m.list.SelectedItem().(Item)
			// if ok {
			// 	m.output.update(i.Title())
			// }
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
	height := (len(*userOptions) * 2) + 4

	m := model{
		list:   list.New(*userOptions, itemDelegate{}, listWidth, height),
		output: output,
	}
	m.list.Title = *title
	m.list.SetShowStatusBar(false)
	m.list.Styles.Title = theme.DefaultStyles().ListTitle
	m.list.Styles.TitleBar = theme.DefaultStyles().ListTitleBar
	return &m
}
