package logger

import "log"

const (
	WarningLogger = "WARNING"
	InfoLogger    = "INFO"
	ErrorLogger   = "ERROR"
)

func Info(a ...interface{}) {
	for _, e := range a {
		log.Printf("[%s]: %+s", InfoLogger, e)
	}
}

func Error(a ...interface{}) {
	for _, e := range a {
		log.Printf("[%s]: %+s", ErrorLogger, e)
	}
}

func Warn(a ...interface{}) {
	for _, e := range a {
		log.Printf("[%s]: %+s", WarningLogger, e)
	}
}
