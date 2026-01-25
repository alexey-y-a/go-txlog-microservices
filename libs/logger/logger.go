package logger

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
    log zerolog.Logger
    once sync.Once
)

func Init() {
	once.Do(func() {
		zerolog.TimeFieldFormat = time.RFC3339
		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		log = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger()
	})
}

func L() zerolog.Logger {
    Init()
    return log
}