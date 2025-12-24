package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init() {
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
}

func GetLogger() zerolog.Logger {
	return log
}

func SetLogger(newLogger zerolog.Logger) {
	log = newLogger
}
