package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var Log zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}
