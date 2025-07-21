package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(logLevel string) {
	level := zerolog.InfoLevel

	if logLevel == "debug" {
		level = zerolog.DebugLevel
	}

	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05"}
	zerolog.SetGlobalLevel(level)
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}

func Get() zerolog.Logger {
	return log.Logger
}
