package logger

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func Init() {
	InfoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	WarningLogger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime)
}

func Info(format string, a ...interface{}) {
	InfoLogger.Printf(format, a...)
}

func Error(format string, a ...interface{}) {
	ErrorLogger.Printf(format, a...)
}

func Warn(format string, a ...interface{}) {
	WarningLogger.Printf(format, a...)
}
