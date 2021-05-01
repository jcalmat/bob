package ui

import (
	ui "github.com/jcalmat/termui/v3"
	"github.com/jcalmat/termui/v3/widgets"
)

type Header struct {
	Logo *widgets.Paragraph
	Help *widgets.Paragraph
}

func NewHeader() *Header {
	h := &Header{}
	h.Logo = h.buildHeader()
	h.Help = h.buildHelp()
	return h
}

func (h *Header) Render() {
	ui.Render(h.Logo, h.Help)
}

func (h *Header) Resize() {
	x, _ := ui.TerminalDimensions()

	h.Logo.SetRect(0, 0, x/2, 12)
	h.Help.SetRect(x, 0, x/2, 12)
}

func (h *Header) HandleEvent(e ui.Event) {}

func (h *Header) buildHeader() *widgets.Paragraph {
	header := widgets.NewParagraph()

	header.Text = `
  ______     ______     ______
 /\  == \   /\  __ \   /\  == \
 \ \  __<   \ \ \/\ \  \ \  __<
  \ \_____\  \ \_____\  \ \_____\
   \/_____/   \/_____/   \/_____/

 Language agnostic Boilerplate Builder
 `
	return header
}

func (h *Header) buildHelp() *widgets.Paragraph {
	help := widgets.NewParagraph()

	help.Title = "Help"
	help.Text = `
    -----------------------------
    -        Move around        -
    -----------------------------
    go up                 ▲
    go down               ▼
    go back               'escape'
    select                'enter'
    quit without saving   'ctrl+c'
`

	return help
}
