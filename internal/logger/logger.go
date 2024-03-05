package logger

import "log"

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// Logger represents a custom logger with different log levels.
type Logger struct {
	level LogLevel
}

// NewLogger creates a new Logger with the specified log level.
func NewLogger(isDebug bool) *Logger {
	if isDebug {
		return &Logger{level: DebugLevel}
	}
	return &Logger{level: InfoLevel}
}

// Debug logs messages at the Debug level.
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DebugLevel {
		format = "[DEBUG] " + format
		log.Printf(format, v...)
	}
}

// Info logs messages at the Info level.
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= InfoLevel {
		format = "[INFO] " + format
		log.Printf(format, v...)
	}
}

// Warn logs messages at the Warn level.
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WarnLevel {
		format = "[WARN] " + format
		log.Printf(format, v...)
	}
}

// Error logs messages at the Error level.
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ErrorLevel {
		format = "[ERROR] " + format
		log.Printf(format, v...)
	}
}
