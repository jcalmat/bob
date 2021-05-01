package command

import "github.com/jcalmat/bob/cmd/cli/ui"

func (c Command) Init(_ ...string) {
	err := c.ConfigApp.InitConfig()
	if err != nil {
		modale := ui.NewModale(err.Error(), ui.ModaleTypeErr)
		modale.Resize()
		modale.Render()
		c.Screen.SetModale(modale)
		return
	}
}
