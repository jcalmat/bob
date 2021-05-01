package ui

import ui "github.com/jcalmat/termui/v3"

type Item interface {
	HandleEvent(e ui.Event)
	Render()
	Resize()
}
