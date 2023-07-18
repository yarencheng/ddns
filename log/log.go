package log

import (
	stdlog "log"
	"os"
	"strings"
	"time"

	zero "github.com/rs/zerolog"
	zerolog "github.com/rs/zerolog/log"
)

var (
	LOG_LEVEL  = strings.ToLower(os.Getenv("LOG_LEVEL"))
	LOG_FORMAT = strings.ToLower(os.Getenv("LOG_FORMAT"))
)

func init() {

	initialize()
}

func initialize() {

	if LOG_LEVEL == "" {
		LOG_LEVEL = "info"
	}

	if LOG_FORMAT == "" {
		LOG_FORMAT = "console"
	}

	zerolog.Info().
		Str("LOG_LEVEL", LOG_LEVEL).
		Str("LOG_FORMAT", LOG_FORMAT).
		Msgf("init logger")

	l := zerolog.With().Caller().Logger()

	switch LOG_FORMAT {
	case "json":
		zero.TimeFieldFormat = zero.TimeFormatUnixMs
		l = l.Output(os.Stdout)
	case "console":
		l = l.Output(
			zero.ConsoleWriter{
				Out:        os.Stdout,
				NoColor:    false,
				TimeFormat: time.RFC3339,
			},
		)
	}

	if level, err := zero.ParseLevel(LOG_LEVEL); err != nil {
		zerolog.Warn().Err(err).Msgf("LOG_LEVEL is invalid. Set to info")
		l = l.Level(zero.DebugLevel)
	} else {
		l = l.Level(level)
	}

	zerolog.Logger = l

	stdlog.SetFlags(0)
	stdlog.SetOutput(zerolog.Logger)
}
