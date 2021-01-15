package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func New(debug bool) zerolog.Logger {
	// set zerolog log level to info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.FormatLevel = func(i interface{}) string {
		return strings.Title(fmt.Sprintf("%s:", i))
	}
	output.FormatTimestamp = func(i interface{}) string { return "" }
	return zerolog.New(output).With().Logger()
}
