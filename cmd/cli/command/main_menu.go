package command

import (
	"errors"

	"github.com/jcalmat/bob/cmd/cli/ui"
	"github.com/jcalmat/bob/pkg/config"
)

func (c Command) MainMenu(args ...string) {

	mainMenu := ui.NewMenu()

	_, err := c.ConfigApp.Parse()
	// if there is no config file, directly go to the init menu
	if err != nil && errors.Is(err, config.ErrConfigFileNotFound) {
		mainMenu.AddOption(ui.MenuOption{
			Name:    "Init",
			Handler: c.Init,
			Description: `Hi, we've detected that you don't have a configuration file for Bob yet, but don't worry, we've got your back!
	
				In order to use Bob, you need to initialize a configuration file.
				This file will contain all the information needed to build your projects. 

				You will be able to specify the templates you want to use and to customize them on the fly to create your own modular templating system.

				This configuration file will be stored in your home directory, in a file called .bobconfig.[yaml/json] depending on the syntax you prefer.

				For your first time, we'll help you create one with some examples already included.

				Thanks for using Bob, we hope you'll enjoy it!`,
		})
	} else {
		mainMenu.AddOptions([]ui.MenuOption{
			{
				Name:        "Build",
				Description: "Build your project",
				Handler:     c.BuildMenu,
			},
			{
				Name:        "Init",
				Handler:     c.Init,
				Description: "Initialize a new configuration file (.bobconfig)",
			},
		})
	}
	mainMenu.AddOptions([]ui.MenuOption{
		{
			Name:        "Help",
			Description: "Help contains some explanations about Bob, indicating what it is precisely and how to use it.",
			Handler:     c.Help,
		},
	})

	mainMenu.Build()

	c.Screen.SetMenu(mainMenu)
	c.Screen.Render()
}
