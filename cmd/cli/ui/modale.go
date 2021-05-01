package ui

import (
	ui "github.com/jcalmat/termui/v3"
	"github.com/jcalmat/termui/v3/widgets"
)

type ModaleOption struct {
	Name        string
	Description string
	Handler     func(...string)
}

type ModaleType int

const (
	ModaleTypeInfo ModaleType = iota
	ModaleTypeErr
	ModaleTypeWrn
)

type Modale struct {
	Text *widgets.Paragraph
	Type ModaleType
}

func NewModale(value string, t ModaleType) *Modale {
	p := widgets.NewParagraph()
	p.Text = value
	return &Modale{
		Text: p,
		Type: t,
	}
}

func (m *Modale) Render() {
	ui.Render(m.Text)
}

func (m *Modale) Resize() {
	x, y := ui.TerminalDimensions()
	_ = x
	_ = y
	m.Text.SetRect(x/2+int(0.2*float64(x/2)), y/2+int(0.2*float64(y/2)), x/2-int(0.2*float64(x/2)), y/2-int(0.2*float64(y/2)))
}

func (m *Modale) HandleEvent(e ui.Event) {
	// switch e.ID {
	// case "<Enter>":
	// 	fn := m.menuOptions[m.Options.Rows[m.Options.SelectedRow][2:]].Handler
	// 	if fn != nil {
	// 		fn(m.Options.Rows[m.Options.SelectedRow][2:])
	// 	}
	// }
	m.Render()
}
