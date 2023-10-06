package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Get() *zerolog.Logger {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &logger
}
