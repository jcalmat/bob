package ui

import (
	ui "github.com/jcalmat/termui/v3"
	"github.com/jcalmat/termui/v3/widgets"
)

type MenuOption struct {
	Name        string
	Description string
	Handler     func(...string)
}

type Menu struct {
	Options     *widgets.List
	Description *widgets.Paragraph

	menuOptions map[string]MenuOption
}

func NewMenu() *Menu {
	return &Menu{
		menuOptions: make(map[string]MenuOption),
	}
}

func (m *Menu) AddOption(o MenuOption) {
	m.menuOptions[o.Name] = o
}

func (m *Menu) AddOptions(os []MenuOption) {
	for _, o := range os {
		m.menuOptions[o.Name] = o
	}
	m.buildOptions(os)
	m.buildDescription()
}

func (m *Menu) Render() {
	ui.Render(m.Options)
	m.Description.Text = m.menuOptions[m.Options.Rows[m.Options.SelectedRow]].Description
	ui.Render(m.Description)
}

func (m *Menu) Resize() {
	x, y := ui.TerminalDimensions()

	m.Options.SetRect(0, 12, x/2, y)
	m.Description.SetRect(x, 12, x/2, y)
}

func (m *Menu) HandleEvent(e ui.Event) {
	switch e.ID {
	case "<Down>":
		m.Options.ScrollDown()
	case "<Up>":
		m.Options.ScrollUp()
	case "<Enter>":
		fn := m.menuOptions[m.Options.Rows[m.Options.SelectedRow][2:]].Handler
		if fn != nil {
			fn(m.Options.Rows[m.Options.SelectedRow])
		}
	}
}

func (m *Menu) buildOptions(os []MenuOption) {
	list := widgets.NewList()
	list.Title = "Main menu"
	list.SelectedRowStyle = ui.NewStyle(ui.ColorClear)
	list.Rows = make([]string, 0)

	for _, o := range os {
		list.Rows = append(list.Rows, "- "+o.Name)
	}

	m.Options = list
}

func (m *Menu) buildDescription() {
	desc := widgets.NewParagraph()
	desc.Title = "Description"
	m.Description = desc
}
