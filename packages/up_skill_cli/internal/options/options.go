package options

import "strings"

type Option string

const (
	Markdown Option = "markdown"
	Man      Option = "man"
	Rest     Option = "rest"
	Yes      Option = "Yes"
	No       Option = "No"
)

var AllowedDocFormats = []string{
	string(Markdown),
	string(Man),
	string(Rest),
}

func AllowedDocFormatsStr() string {
	return strings.Join(AllowedDocFormats, ", ")
}

type Item struct {
	title Option
	desc  string
}

func NewItem(title Option, desc string) Item {
	return Item{title: title, desc: desc}
}

func (i Item) Title() string       { return string(i.title) }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return string(i.title) }
