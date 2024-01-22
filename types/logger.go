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

func (l *Logger) Init() {
	l.infoLogger = *log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	l.debugLogger = *log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime)
	l.warningLogger = *log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime)
	l.errorLogger = *log.New(os.Stdout, "[Error] ", log.Ldate|log.Ltime)
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

func (l *Logger) Fatal(format string, v ...any) {
	l.errorLogger.Fatalf(format, v...)
}
