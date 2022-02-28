package internal

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func init() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339Nano,
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano
	Logger = zerolog.
		New(output).
		With().
		Timestamp().
		Logger()
}

func SetLogLevel(s string) error {
	l, err := zerolog.ParseLevel(s)
	if err != nil {
		return err
	}

	Logger = Logger.Level(l)
	return nil
}
