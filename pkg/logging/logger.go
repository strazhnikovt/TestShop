package logging

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[APP] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Printf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(format, v...)
}
