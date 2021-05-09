package command

import "github.com/jcalmat/bob/cmd/cli/ui"

func (c Command) Init(_ ...string) {
	//TODO: create menu to choose type of config file
	err := c.ConfigApp.InitConfig()
	if err != nil {
		c.Screen.RenderModale(err.Error(), ui.ModaleTypeErr)
		return
	}
}
