package unwindlogger

type Level int

const (
	TRACE    Level = 0
	DEBUG    Level = 20
	INFO     Level = 40
	WARN     Level = 60
	ERROR    Level = 80
	CRITICAL Level = 100
)

type Logger struct {
	level Level
}

func NewLogger() *Logger {
	return &Logger{
		level: INFO,
	}
}

func (l *Logger) WithLevel(level Level) *Logger {
	l.level = level
	return l
}

func (l *Logger) WithField(key string, value interface{}) *Entry {
	return NewEntry(l).WithField(key, value)
}

func (l *Logger) Debug(msg string) {
	NewEntry(l).Debug(msg)
}
