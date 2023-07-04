package unwindlogger_test

import (
	"testing"
	"unwindlogger"
)

func TestLogger(t *testing.T) {
	l := unwindlogger.NewLogger()

	l.WithField("hello", "world").Debug("the message!")
}
