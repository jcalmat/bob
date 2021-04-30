package main

import (
	"log"

	"github.com/jcalmat/bob/cmd/cli/command"
	"github.com/jcalmat/bob/cmd/cli/ui"
	"github.com/jcalmat/bob/pkg/config/app"
	"github.com/jcalmat/bob/pkg/logger"
	"github.com/mitchellh/go-homedir"
)

var (
	configFilePath string = "~/.bobconfig.yml"
)

func main() {
	// io.ASCIIBob()

	logger := logger.New(false)

	absPath, _ := homedir.Expand(configFilePath)
	configApp := app.App{
		ConfigFilePath: absPath,
	}

	//

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	screen := ui.NewScreen()
	// defer ui.Close()

	handler := command.Command{
		Logger: logger,

		ConfigApp: configApp,
		Screen:    screen,
	}

	mainMenu := ui.NewMenu()
	mainMenu.AddOptions([]ui.MenuOption{
		{
			Name:    "Build",
			Handler: handler.BuildMenu,
			Description: `
			Build a project from a specified template
			`,
		},
		{
			Name:    "Init",
			Handler: handler.Init,
			Description: `
			Initialize a new .bobconfig file if it doesn't already exist.
			`,
		},
		{
			Name: "Who am I?",
			Description: `
				I am a tool used to generate boilerplate code from templates.

				I use go templates syntaxe to parse and replace variables, thus these variables must be formatted with double brackets like {{VARIABLE}}.
				For more information about the format, here is a cheat sheet: https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet.
				
			`,
		},
	})

	screen.SetMenu(mainMenu)

	screen.Render()
	screen.HandleEvents()
	//

	// r := cli.Init(logger)
	// buildCmd := r.Command.AddCommand(&cli.Command{
	// 	Key:         "build",
	// 	Description: "build a project from a specified template",
	// })

	// globalConfig, err := configApp.Parse()
	// if err != nil {
	// 	logger.Err(err).Msg("")
	// 	return
	// }
	// for key := range globalConfig.Commands {
	// 	buildCmd.AddCommand(&cli.Command{
	// 		Key:     key,
	// 		Handler: handler.Build,
	// 	})
	// }

	// r.Command.AddCommand(&cli.Command{
	// 	Key:         "init",
	// 	Description: "initialize bob's config file",
	// 	Handler:     handler.Init,
	// })

	// r.Handle()
}
