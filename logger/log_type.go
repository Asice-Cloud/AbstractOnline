package logger

import (
	"time"
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warn
	Panic
	ERROR
	Fatal
)

type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
}

func (le LogEntry) String() string {
	return le.Timestamp.Format(time.RFC3339) + " [" + le.levelString() + "] " + le.Message
}

func (le LogEntry) levelString() string {
	switch le.Level {
	case Debug:
		return "Debug"
	case Info:
		return "Info"
	case Warn:
		return "Warn"
	case Panic:
		return "Panic"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
