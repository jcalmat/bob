package ui

import (
	ui "github.com/jcalmat/termui/v3"
)

type Screen struct {
	Headers *Header
	Menu    *Menu
	Form    *Form
}

func NewScreen() *Screen {
	return &Screen{
		Headers: NewHeader(),
	}
}

func (s *Screen) SetMenu(m *Menu) {
	s.Menu = m
	s.UnsetForm()
}

func (s *Screen) SetForm(f *Form) {
	s.Form = f
	s.UnsetMenu()
}

func (s *Screen) UnsetMenu() {
	s.Menu = nil
}

func (s *Screen) UnsetForm() {
	s.Form = nil
}

func Init() error {
	return ui.Init()
}

func Close() {
	ui.Close()
}

func (s *Screen) Render() {
	s.Resize()
	if s.Headers != nil {
		s.Headers.Render()
	}

	if s.Menu != nil {
		s.Menu.Render()
	}

	if s.Form != nil {
		s.Menu.Render()
	}
}

func (s *Screen) Resize() {
	if s.Headers != nil {
		s.Headers.Resize()
	}

	if s.Menu != nil {
		s.Menu.Resize()
	}

	if s.Form != nil {
		s.Menu.Resize()
	}
}

func (s *Screen) HandleEvents() {
	var close bool
	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		if s.Menu != nil {
			s.Menu.HandleEvent(e)
		}

		if s.Form != nil {
			s.Menu.HandleEvent(e)
		}
		switch e.ID {
		case "<C-c>":
			ui.Close()
			close = true
		case "<Resize>":
			s.Resize()
		}

		// escape = restore prev screen state

		if close {
			break
		}
		s.Render()
	}
}
