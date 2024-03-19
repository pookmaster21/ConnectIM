package types

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger    log.Logger
	debugLogger   log.Logger
	warningLogger log.Logger
	errorLogger   log.Logger
}

var logger *Logger

func NewLogger() *Logger {
	if logger == nil {
		logger = new(Logger)
		logger.infoLogger = *log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
		logger.debugLogger = *log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime)
		logger.warningLogger = *log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime)
		logger.errorLogger = *log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime)
	}

	return logger
}

func (l *Logger) Info(format string, v ...any) {
	l.infoLogger.Printf(format, v...)
}

func (l *Logger) Debug(format string, v ...any) {
	l.debugLogger.Printf(format, v...)
}

func (l *Logger) Warn(format string, v ...any) {
	l.warningLogger.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...any) {
	l.errorLogger.Printf(format, v...)
}
