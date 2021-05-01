package command

import "github.com/jcalmat/bob/cmd/cli/ui"

func (c Command) BuildMenu(args ...string) {
	globalConfig, err := c.ConfigApp.Parse()
	if err != nil {
		c.Logger.Err(err).Msg("")
		return
	}

	menu := ui.NewMenu()

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
