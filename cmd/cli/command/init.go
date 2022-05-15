package command

import (
	"github.com/jcalmat/bob/cmd/cli/ui"
	"github.com/jcalmat/bob/pkg/config"
	"github.com/jcalmat/termui/v3/widgets"
	"github.com/mitchellh/go-homedir"
)

func (c Command) Init(_ ...string) {

	menu := ui.NewMenu()
	menu.Options.Title = "Which syntax do you prefer?"

	menu.AddOptions([]ui.MenuOption{
		{
			Name:        "json",
			Handler:     c.initConfig,
			Description: "Initialize a new .bobconfig with json format",
		},
		{
			Name:        "yaml",
			Handler:     c.initConfig,
			Description: "Initialize a new .bobconfig with yaml format",
		},
	})

	menu.Build()

	c.Screen.SetMenu(menu)
	c.Screen.Render()
}

func (c Command) initConfig(types ...string) {
	t := types[0]

	config := config.C{
		Commands: map[string]config.Command{
			"json": {
				Git: "https://github.com/jcalmat/bob_templates json_config",
			},
			"yaml": {
				Git: "https://github.com/jcalmat/bob_templates yaml_config",
			},
		},
	}

	absPath, _ := homedir.Expand("~")

	builder := &Builder{
		form:            ui.NewForm(),
		screen:          c.Screen,
		settings:        config.Settings,
		skipMap:         make(map[string]struct{}),
		stringQuestions: make(map[string]*widgets.TextField),
		boolQuestions:   make(map[string]*widgets.Checkbox),
		configApp:       c.ConfigApp,
		command:         config.Commands[t],
		forcedPath:      &absPath,
	}

	builder.SetupBuild(c.Screen)
}
