package unwindlogger_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/Ben435/unwindlogger"
)

func TestLogger(t *testing.T) {
	buffer := &bytes.Buffer{}

	l := unwindlogger.NewLogger().
		WithLevel(unwindlogger.DEBUG).
		WithOut(buffer)

	l.WithField("hello", "world").Debug("debug message!")
	l.Info("info message!")
	l.Warn("warn message!")
	l.WithError(fmt.Errorf("bad error")).Error("error message!")

	allWrites := buffer.String()
	assert.Contains(t, allWrites, "debug message!")
	assert.Contains(t, allWrites, "\"hello\":\"world\"")

	assert.Contains(t, allWrites, "info message!")

	assert.Contains(t, allWrites, "warn message!")

	assert.Contains(t, allWrites, "error message!")
	assert.Contains(t, allWrites, "\"error\":\"bad error\"")
}
