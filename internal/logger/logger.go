// logger handles setting the default output format for logging
package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// var Logger zerolog.Logger //nolint:gochecknoglobals // allow for logging at this time.

// InitLogger initializes the zerolog logger settings.
func New() zerolog.Logger {
	l := zerolog.New(os.Stdout).With().Logger().With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	return l
}

// EnableDebug enables debug logging.
func EnableDebug() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
