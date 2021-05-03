package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jcalmat/bob/cmd/cli/command"
	"github.com/jcalmat/bob/cmd/cli/ui"
	"github.com/jcalmat/bob/pkg/config/app"
	"github.com/mitchellh/go-homedir"
)

var (
	configFilePath string = "~/.bobconfig.yml"
	version        string
)

func main() {

	// handle flags
	displayVersion := flag.Bool("v", false, "display bob version")
	flag.Parse()
	if *displayVersion {
		fmt.Println(version)
		return
	}

	absPath, _ := homedir.Expand(configFilePath)
	configApp := app.App{
		ConfigFilePath: absPath,
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	screen := ui.NewScreen()

	handler := command.Command{
		ConfigApp: configApp,
		Screen:    screen,
	}

	mainMenu := ui.NewMenu()
	mainMenu.AddOptions([]ui.MenuOption{
		{
			Name:        "Build",
			Handler:     handler.BuildMenu,
			Description: "Build a project from a specified template",
		},
		{
			Name:        "Init",
			Handler:     handler.Init,
			Description: "Initialize a new .bobconfig file if it doesn't already exist.",
		},
		{
			Name:        "Help",
			Description: "Help contains some explanations about Bob, indicating what it is precisely and how to use it.",
			Handler:     handler.Help,
		},
	})

	screen.SetMenu(mainMenu)

	screen.Run()
}
