package unwindlogger

import (
	"io"
	"os"
)

type Logger struct {
	level Level
	out   io.Writer
}

func NewLogger() *Logger {
	return &Logger{
		level: INFO,
		out:   os.Stdout,
	}
}

func (l *Logger) WithLevel(level Level) *Logger {
	l.level = level
	return l
}

func (l *Logger) WithField(key string, value interface{}) *Entry {
	return NewEntry(l).WithField(key, value)
}

func (l *Logger) WithError(err error) *Entry {
	return NewEntry(l).WithError(err)
}

func (l *Logger) Debug(msg string) {
	NewEntry(l).Debug(msg)
}

func (l *Logger) Info(msg string) {
	NewEntry(l).Info(msg)
}

func (l *Logger) Warn(msg string) {
	NewEntry(l).Warn(msg)
}

func (l *Logger) Error(msg string) {
	NewEntry(l).Error(msg)
}
