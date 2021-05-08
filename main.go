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

	handler.MainMenu()

	screen.Run()
}
