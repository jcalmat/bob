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

func NewModale(s string, t ModaleType) *Modale {
	p := widgets.NewParagraph()
	p.Text = s
	return &Modale{
		Text: p,
		Type: t,
	}
}

// Render implements Item interface
func (m *Modale) Render() {
	ui.Render(m.Text)
}

// Resize implements Item interface
func (m *Modale) Resize() {
	x, y := ui.TerminalDimensions()
	_ = x
	_ = y
	m.Text.SetRect(x/2+int(0.2*float64(x/2)), y/2+int(0.2*float64(y/2)), x/2-int(0.2*float64(x/2)), y/2-int(0.2*float64(y/2)))
}

// HandleEvent implements Item interface
func (m *Modale) HandleEvent(e ui.Event) {}

func (s *Screen) RenderModale(msg string, t ModaleType) {
	modale := NewModale(msg, t)
	modale.Resize()
	modale.Render()
	s.SetModale(modale)
}
