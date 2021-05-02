package command

import "github.com/jcalmat/bob/cmd/cli/ui"

func (c Command) Help(args ...string) {

	menu := ui.NewMenu()

	menu.AddOptions([]ui.MenuOption{
		{
			Name: "What is Bob?",
			Description: `
			Bob is a tool used to generate boilerplate code.

			It was made to avoid loosing too much time writing redundant code and allow developers to focus and dedicate more on interesting parts of their projects.

			Bob is asking you to do 2 things in order to work:
			- Create template(s) of boilerplate code
			- Registering those templates to a config file written in yaml or json in your root directory. This file shall be named .bobconfig.yaml/.bobconfig.yml/.bobconfig.json
		`,
		},
		{
			Name: "How does it work?",
			Description: `
			## Introduction

			Bob is a boilerplate generator. To generate this boilerplate in your project, you need to create templates which represent roughly what your code should look like in the end.

			Though, even if boilerplate represents redundant code, it could be slightly different from one file to another. Variables can change, blocks of code can exist in one but not in the other, this kind of things.
			This is why Bob addresses this issue by using a particular syntax to perform some modifications to your boiler. 
			
			## Templates syntax

			Bob uses go templates syntax to parse and replace variables, here is a quick overview of go template documentation:
			
			"The input text for a template is UTF-8-encoded text in any format.
			"Actions"--data evaluations or control structures--are delimited by "{{" and "}}"; all text outside actions is copied to the output unchanged. Except for raw strings, actions may not span newlines, although comments can."


			For more information about the format, here is a cheat sheet: https://golang.org/pkg/text/template/#hdr-Actions.

			`,
		},
	})

	menu.Options.Title = "Help"
	menu.Description.Title = ""
	c.Screen.SetMenu(menu)
	menu.Render()

	// globalConfig, err := c.ConfigApp.Parse()
	// if err != nil {
	// 	modale := ui.NewModale(err.Error(), ui.ModaleTypeErr)
	// 	modale.Render()
	// 	c.Screen.SetModale(modale)
	// 	return
	// }

	// options := make([]ui.MenuOption, 0)
	// for key, command := range globalConfig.Commands {
	// 	options = append(options, ui.MenuOption{
	// 		Name:        key,
	// 		Description: command.Description,
	// 		Handler:     c.Build,
	// 	})
	// }

	// menu.AddOptions(options)

	// menu.Options.Title = "Build menu"
	// c.Screen.SetMenu(menu)

	// menu.Render()
}
