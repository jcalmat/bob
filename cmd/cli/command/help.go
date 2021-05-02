package command

import "github.com/jcalmat/bob/cmd/cli/ui"

func (c Command) Help(args ...string) {

	menu := ui.NewMenu()

	menu.AddOptions([]ui.MenuOption{
		{
			Name: "What is Bob?",
			Description: `## Bob
			
			Bob, or Boilerplate Builder, is a tool used to generate boilerplate code.

			It was made to avoid loosing too much time writing redundant code and allow developers to focus and dedicate more time on interesting parts of their projects.

			Bob is asking you to do 2 things in order to work:
			- Create template(s) of boilerplate code
			- Registering those templates to a config file written in yaml or json in your root directory. This file shall be named .bobconfig.yaml/.bobconfig.yml/.bobconfig.json
		`,
		},
		{
			Name: "How does it work?",
			Description: `## Introduction

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
		{
			Name: "Examples - The basics",
			Description: `## The Basics
			
			Given the file "example.js" with the following line of code:

			var {{.my_variable}}
			
			Given the following .bobconfig.yaml:

			commands:
			  test:
			    templates:
			      - test
			
			templates:
			  test:
			    path: "/path/to/example.js"
			    variables:
			      - name: "my_variable"
			        type: "string"
			
			When you run bob, it will automatically propose you to replace "my_variable" by any string and ask you where to put your newly created boiler file.
			`,
		},
		{
			Name: "Examples - Special variables",
			Description: `## Special variables
			
			Since bob uses go template to perform the variable replacement, it has some interesting specificities.
			You can, for instance, perform conditional operations.

			Ex:

			{{if .my_variable}}
			// do something only if .my_variable is defined
			{{end}}

			Bob also ships with homemade functions to ease string formatting
			- **short [INT]** will truncate the x first characters of your variable
			- **upcase** will capitalize your variable
			- **title** will return a copy of the string s with all Unicode letters that begin words mapped to their Unicode title case

			Ex:

			{{title .my_variable}}
			// .my_variable = test -> Test

			{{short .my_variable 1}}
			// .my_variable = test -> t

			{{upcase .my_variable}}
			// .my_variable = test -> TEST


			You can even combine multiple functions.

			Ex:

			{{short .my_variable 3 | upcase}}
			// my_variable = test -> TES
			`,
		},
	})

	menu.Options.Title = "Help"
	menu.Description.Title = ""
	c.Screen.SetMenu(menu)
	menu.Render()
}
