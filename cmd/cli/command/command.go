package command

import (
	"github.com/jcalmat/bob/cmd/cli/ui"
	"github.com/jcalmat/bob/pkg/config"
)

type Command struct {
	ConfigApp config.App
	Screen    *ui.Screen
}
