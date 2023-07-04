package unwindlogger_test

import (
	"fmt"
	"testing"

	"github.com/Ben435/unwindlogger"
)

func TestLogger(t *testing.T) {
	l := unwindlogger.NewLogger().WithLevel(unwindlogger.DEBUG)

	l.WithField("hello", "world").Debug("debug message!")
	l.Info("info message!")
	l.Warn("the message!")
	l.WithError(fmt.Errorf("bad error")).Error("error message!")
}
