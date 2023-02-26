package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger zerolog.Logger
)

func Info(format string, v ...interface{}) {
	logger.Info().CallerSkipFrame(1).Msg(fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	logger.Error().CallerSkipFrame(1).Msg(fmt.Sprintf(format, v...))
}

func Warn(format string, v ...interface{}) {
	logger.Warn().CallerSkipFrame(1).Msg(fmt.Sprintf(format, v...))
}

func Init(debug bool) {
	fmt.Println("Init logger..")
	var writers []io.Writer
	if debug {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		writers = append(writers, os.Stdout)
	}
	writers = append(writers, newRollingFile())
	mw := io.MultiWriter(writers...)
	logger = zerolog.New(mw).With().Timestamp().Caller().Logger()
}

func newRollingFile() io.Writer {
	pathfile := path.Join("/tmp/product-crud", "log.json")
	if _, err := os.Stat(pathfile); err != nil {
		fmt.Println("Path does not exist!", err)
	}
	return &lumberjack.Logger{
		Filename:   path.Join(pathfile),
		MaxBackups: 5,
		MaxSize:    50,
	}
}
