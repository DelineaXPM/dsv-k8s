// logger handles setting the default output format for logging
package logger

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
)

// var Logger zerolog.Logger //nolint:gochecknoglobals // allow for logging at this time.

// InitLogger initializes the zerolog logger settings.
func New() zerolog.Logger {
	// see zerolog docs, this removes the fully qualified path and makes it just show the "logger.go:linenumber" so logs are much easier to read
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	l := zerolog.New(os.Stdout).With().Caller().Logger().With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	return l
}

// EnableDebug enables debug logging.
func EnableDebug() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
