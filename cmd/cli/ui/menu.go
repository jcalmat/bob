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

	menuOptions    map[string]MenuOption
	orderedOptions []MenuOption
}

func NewMenu() *Menu {
	return &Menu{
		Options:     &widgets.List{},
		Description: &widgets.Paragraph{},
		menuOptions: make(map[string]MenuOption),
	}
}

func (m *Menu) AddOption(o MenuOption) {
	m.menuOptions[o.Name] = o
	m.orderedOptions = append(m.orderedOptions, o)
}

func (m *Menu) AddOptions(os []MenuOption) {
	for _, o := range os {
		m.AddOption(o)
	}
}

func (m *Menu) Build() {
	m.buildOptions(m.orderedOptions)
	m.buildDescription()
}

// Render implements Item interface
func (m *Menu) Render() {
	ui.Render(m.Options)
	if len(m.Options.Rows) > m.Options.SelectedRow {
		if opt, ok := m.menuOptions[m.Options.Rows[m.Options.SelectedRow][2:]]; ok {
			m.Description.Text = opt.Description
		}
	}
	ui.Render(m.Description)
}

// Resize implements Item interface
func (m *Menu) Resize() {
	x, y := ui.TerminalDimensions()

	m.Options.SetRect(0, 12, x/2, y)
	m.Description.SetRect(x, 12, x/2, y)
}

// HandleEvent implements Item interface
func (m *Menu) HandleEvent(e ui.Event) {
	switch e.ID {
	case "<Down>":
		m.Options.ScrollDown()
	case "<Up>":
		m.Options.ScrollUp()
	case "<Enter>":
		fn := m.menuOptions[m.Options.Rows[m.Options.SelectedRow][2:]].Handler
		if fn != nil {
			fn(m.Options.Rows[m.Options.SelectedRow][2:])
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
