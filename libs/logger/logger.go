package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var (
    log zerolog.Logger
    inited bool
)

func Init() {
    log = zerolog.New(os.Stdout).With().Timestamp().Logger()
    inited = true
}

func L() zerolog.Logger {
    if !inited {
        Init()
    }
    return log
}