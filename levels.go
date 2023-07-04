package unwindlogger

import "strconv"

type Level int

const (
	DEBUG Level = 20
	INFO  Level = 40
	WARN  Level = 60
	ERROR Level = 80
)

func (l Level) Format() string {
	switch l {
	case DEBUG:
		{
			return "DEBUG"
		}
	case INFO:
		{
			return "INFO"
		}
	case WARN:
		{
			return "WARN"
		}
	case ERROR:
		{
			return "ERROR"
		}
	default:
		{
			return strconv.Itoa(int(l))
		}
	}
}
