package ui

import (
	ui "github.com/jcalmat/termui/v3"
)

const (
	StatusRunning = iota
	StatusStopped
)

type Screen struct {
	Headers *Header
	Menu    *Menu
	Form    *Form

	status     int
	breadcrumb []Screen
}

func NewScreen() *Screen {
	return &Screen{
		Headers:    NewHeader(),
		breadcrumb: make([]Screen, 0),
	}
}

func (s *Screen) SetMenu(m *Menu) {
	s.breadcrumb = append(s.breadcrumb, *s)
	s.Menu = m
	s.Menu.Resize()
	s.UnsetForm()
}

func (s *Screen) SetForm(f *Form) {
	s.breadcrumb = append(s.breadcrumb, *s)
	s.Form = f
	s.Form.Resize()
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

func (s *Screen) Run() {
	s.status = StatusRunning
	s.Resize()
	s.Render()
	s.HandleEvents()
}

func (s *Screen) Stop() {
	s.status = StatusStopped
	ui.Close()
}

func (s *Screen) Render() {
	if s.Headers != nil {
		s.Headers.Render()
	}

	if s.Menu != nil {
		s.Menu.Render()
	}

	if s.Form != nil {
		s.Form.Render()
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
		s.Form.Resize()
	}
}

func (s *Screen) Restore(old Screen) {
	if old.Headers != nil {
		s.Headers = old.Headers
	}

	if old.Menu != nil {
		s.Menu = old.Menu
	}

	if old.Form != nil {
		s.Form = old.Form
	}

	s.breadcrumb = old.breadcrumb
	s.status = old.status

	s.Render()
}

func (s *Screen) HandleEvents() {
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if s.Form != nil {
			s.Form.HandleEvent(e)
		}

		if s.Menu != nil {
			s.Menu.HandleEvent(e)
		}

		switch e.ID {
		case "<C-c>":
			s.Stop()
		case "<Resize>":
			s.Resize()
		case "<Escape>":
			if len(s.breadcrumb) > 1 {
				s.Restore(s.breadcrumb[len(s.breadcrumb)-1])
			} else {
				s.Stop()
			}
		}

		if s.status == StatusStopped {
			break
		}

		s.Render()
	}
}
