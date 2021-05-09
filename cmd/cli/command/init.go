package command

import (
	"github.com/jcalmat/bob/cmd/cli/ui"
	"github.com/jcalmat/termui/v3"
)

func (c Command) Init(_ ...string) {

	menu := ui.NewMenu()
	menu.Options.Title = "Which syntax do you prefer?"
	menu.AddOptions([]ui.MenuOption{
		{
			Name:        "json",
			Handler:     c.initConfig,
			Description: "Initialize a new .bobconfig with json format",
		},
		{
			Name:        "yaml",
			Handler:     c.initConfig,
			Description: "Initialize a new .bobconfig with yaml format",
		},
	})

	c.Screen.SetMenu(menu)
	c.Screen.Render()
}

func (c Command) initConfig(types ...string) {
	t := types[0]

	switch t {
	case "json":
		err := c.ConfigApp.InitJSONConfig()
		if err != nil {
			c.Screen.RenderModale(err.Error(), ui.ModaleTypeErr)
			return
		}
	case "yaml":
		err := c.ConfigApp.InitYamlConfig()
		if err != nil {
			c.Screen.RenderModale(err.Error(), ui.ModaleTypeErr)
			return
		}
	}

	c.Screen.RenderModale(`
			Done
			Press Enter to quit
			Esc to get back to main menu
		`, ui.ModaleTypeInfo)

	uiEvents := termui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			c.Screen.Stop()
			return
		case "<Enter>":
			c.Screen.Stop()
			return
		case "<Escape>":
			c.Screen.Restore()
			return
		}
	}
}
