package ui

import (
	ui "github.com/jcalmat/termui/v3"
	"github.com/jcalmat/termui/v3/widgets"
)

type Form struct {
	Content *widgets.Form
	Infos   *widgets.Paragraph
}

func NewForm() *Form {
	return &Form{}
}

func (f *Form) SetTitle(title string) {
	if f.Content == nil {
		f.buildContent()
	}
	f.Content.Title = title
}

func (f *Form) SetNodes(n []*widgets.FormNode) {
	if f.Content == nil {
		f.buildContent()
	}
	f.Content.SetNodes(n)
}

func (f *Form) SetInfos(s string) {
	if f.Infos == nil {
		f.buildInfos()
	}
	f.Infos.Text = s
}

func (f *Form) Render() {
	if f.Infos == nil {
		f.buildInfos()
	}

	ui.Render(f.Content)
	ui.Render(f.Infos)
}

func (f *Form) Resize() {
	x, y := ui.TerminalDimensions()

	f.Content.SetRect(0, 12, x/2, y)
	f.Infos.SetRect(x, 12, x/2, y)
}

func (f *Form) HandleEvent(e ui.Event) {
	f.Content.HandleKeyboard(e)

	switch e.ID {
	case "<Down>":
		f.Content.ScrollDown()
	case "<Up>":
		f.Content.ScrollUp()
	case "<Enter>":
		f.Content.ToggleExpand()
		f.Content.ScrollDown()
	case "<Close>":
		break
	}
}

func (f *Form) buildContent() {
	form := widgets.NewForm()
	form.Title = "Main Form"
	form.SelectedTextStyle = ui.NewStyle(ui.ColorClear)
	f.Content = form
}

func (f *Form) buildInfos() {
	desc := widgets.NewParagraph()
	desc.Title = "Infos"
	f.Infos = desc
}
