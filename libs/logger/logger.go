package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init() {
    log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func L() zerolog.Logger {
    return log
}