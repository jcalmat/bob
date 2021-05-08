package command

import "github.com/jcalmat/bob/cmd/cli/ui"

func (c Command) BuildMenu(args ...string) {

	menu := ui.NewMenu()

	globalConfig, err := c.ConfigApp.Parse()
	if err != nil {
		modale := ui.NewModale(err.Error(), ui.ModaleTypeErr)
		modale.Render()
		c.Screen.SetModale(modale)
		return
	}

	options := make([]ui.MenuOption, 0)
	for key, command := range globalConfig.Commands {
		options = append(options, ui.MenuOption{
			Name:        key,
			Description: command.Description,
			Handler:     c.Build,
		})
	}

	menu.AddOptions(options)

	menu.Options.Title = "Build menu"
	c.Screen.SetMenu(menu)

	menu.Render()
}
