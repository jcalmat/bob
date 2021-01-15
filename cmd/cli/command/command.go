package command

import (
	"github.com/jcalmat/bob/pkg/config"

	"github.com/rs/zerolog"
)

type Command struct {
	Logger    zerolog.Logger
	ConfigApp config.App
}
