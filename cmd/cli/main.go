package main

import (
	"github.com/jcalmat/bob/cmd/cli/command"
	"github.com/jcalmat/bob/pkg/cli"
	"github.com/jcalmat/bob/pkg/config/app"
	"github.com/jcalmat/bob/pkg/io"
	"github.com/jcalmat/bob/pkg/logger"
)

var (
	configFilePath string = "/home/kyanpai/.bobconfig.yml"
)

func main() {
	io.ASCIIBob()

	logger := logger.New(false)

	configApp := app.App{
		ConfigFilePath: configFilePath,
	}

	handler := command.Command{
		Logger: logger,

		ConfigApp: configApp,
	}

	r := cli.Init(logger)
	buildCmd := r.Command.AddCommand(&cli.Command{
		Key:         "build",
		Description: "build a project from a specified template",
	})

	globalConfig, err := configApp.Parse()
	if err != nil {
		logger.Err(err).Msg("")
		return
	}
	for key := range globalConfig.Commands {
		buildCmd.AddCommand(&cli.Command{
			Key:     key,
			Handler: handler.Build,
		})
	}

	r.Command.AddCommand(&cli.Command{
		Key:         "init",
		Description: "initialize bob's config file",
		Handler:     handler.Init,
	})

	r.Handle()
}
