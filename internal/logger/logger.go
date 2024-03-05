package logger

import "log"

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Logger struct {
	level LogLevel
}

func NewLogger(isDebug bool) *Logger {
	if isDebug {
		return &Logger{level: DebugLevel}
	}
	return &Logger{level: InfoLevel}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DebugLevel {
		format = "[DEBUG] " + format
		log.Printf(format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= InfoLevel {
		format = "[INFO] " + format
		log.Printf(format, v...)
	}
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WarnLevel {
		format = "[WARN] " + format
		log.Printf(format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ErrorLevel {
		format = "[ERROR] " + format
		log.Printf(format, v...)
	}
}
