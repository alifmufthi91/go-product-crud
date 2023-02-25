package logger

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// var (
// 	Log *zerolog.Logger = &zerolog.Logger{}
// )

func Info(format string, v ...interface{}) {
	log.Info().Msg(fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	log.Error().Msg(fmt.Sprintf(format, v...))
}

func Warn(format string, v ...interface{}) {
	log.Warn().Msg(fmt.Sprintf(format, v...))
}
