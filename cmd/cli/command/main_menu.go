package command

import "github.com/jcalmat/bob/cmd/cli/ui"

func (c Command) MainMenu(args ...string) {

	mainMenu := ui.NewMenu()
	mainMenu.AddOptions([]ui.MenuOption{
		{
			Name:        "Build",
			Handler:     c.BuildMenu,
			Description: "Build a project from a specified template",
		},
		{
			Name:        "Init",
			Handler:     c.Init,
			Description: "Initialize a new .bobconfig file if it doesn't already exist.",
		},
		{
			Name:        "Help",
			Description: "Help contains some explanations about Bob, indicating what it is precisely and how to use it.",
			Handler:     c.Help,
		},
	})

	c.Screen.SetMenu(mainMenu)
	c.Screen.Render()
}
