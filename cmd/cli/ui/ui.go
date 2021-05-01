package ui

import (
	ui "github.com/jcalmat/termui/v3"
)

type Status int

const (
	Running = iota
	Stopped
)

type ItemType int

const (
	HeaderItem ItemType = iota
	CenteredItem
	ModaleItem
)

type Screen struct {
	Headers *Header

	Items       []Item
	FocusedItem Item

	status int
}

func NewScreen() *Screen {
	header := NewHeader()
	header.Resize()
	header.Render()

	return &Screen{
		Items:   make([]Item, 0),
		Headers: header,
	}
}

func (s *Screen) SetMenu(m *Menu) {
	s.Items = append(s.Items, m)
	m.Resize()
	s.FocusedItem = m
}

func (s *Screen) SetForm(f *Form) {
	s.Items = append(s.Items, f)
	f.Resize()
	s.FocusedItem = f
}

func (s *Screen) SetModale(m *Modale) {
	s.Items = append(s.Items, m)
	m.Resize()
	s.FocusedItem = m
}

func Init() error {
	return ui.Init()
}

func Close() {
	ui.Close()
}

func (s *Screen) Run() {
	s.status = Running
	s.Resize()
	s.Render()
	s.HandleEvents()
}

func (s *Screen) Stop() {
	s.status = Stopped
	ui.Close()
}

func (s *Screen) Render() {
	s.FocusedItem.Render()
}

func (s *Screen) Resize() {
	for _, i := range s.Items {
		i.Resize()
	}
}

func (s *Screen) Restore() {
	s.Items = s.Items[:len(s.Items)-1]
	s.FocusedItem = s.Items[len(s.Items)-1]
	s.Render()
}

func (s *Screen) HandleEvents() {
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		s.FocusedItem.HandleEvent(e)

		switch e.ID {
		case "<C-c>":
			s.Stop()
		case "<Resize>":
			s.Resize()
		case "<Escape>":
			if len(s.Items) > 1 {
				s.Restore()
			} else {
				s.Stop()
			}
		}

		if s.status == Stopped {
			break
		}

		s.Render()
	}
}
