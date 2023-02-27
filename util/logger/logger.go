package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"product-crud/config"
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

func Init() {
	log.Println("Init logger..")
	var writers []io.Writer
	if config.GetEnv().Mode != "prod" {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
		writers = append(writers, newRollingFile())
	} else {
		writers = append(writers, os.Stdout)
	}
	mw := io.MultiWriter(writers...)
	logger = zerolog.New(mw).With().Timestamp().Caller().Logger()
}

func newRollingFile() io.Writer {
	pathfile := path.Join(config.GetEnv().FilePath, "product-crud", "log.json")
	if _, err := os.Stat(pathfile); err != nil {
		log.Println("Path does not exist!", err)
	}
	return &lumberjack.Logger{
		Filename:   pathfile,
		MaxBackups: 2,
		MaxSize:    50,
	}
}
