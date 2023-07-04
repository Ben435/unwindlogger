package unwindlogger_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/Ben435/unwindlogger"
)

func TestLogger_DirectLogging(t *testing.T) {
	buffer := &bytes.Buffer{}
	ctx := context.Background()

	l := unwindlogger.NewLogger().
		WithLevel(unwindlogger.DEBUG).
		WithOut(buffer)

	l.WithContext(ctx).WithField("hello", "world").Debug("debug message!")
	l.WithContext(ctx).Info("info message!")
	l.WithContext(ctx).Warn("warn message!")
	l.WithContext(ctx).WithError(fmt.Errorf("bad error")).Error("error message!")

	allWrites := buffer.String()
	assert.Contains(t, allWrites, "debug message!")
	assert.Contains(t, allWrites, "\"hello\":\"world\"")

	assert.Contains(t, allWrites, "info message!")

	assert.Contains(t, allWrites, "warn message!")

	assert.Contains(t, allWrites, "error message!")
	assert.Contains(t, allWrites, "\"error\":\"bad error\"")
}

func TestLogger_UnwindLogging(t *testing.T) {
	buffer := &bytes.Buffer{}
	ctx := context.Background()
	err := fmt.Errorf("bad error")

	l := unwindlogger.NewLogger().
		WithLevel(unwindlogger.WARN).
		WithOut(buffer)

	ctx = l.StartTracking(ctx)

	l.WithContext(ctx).WithField("hello", "world").Debug("debug message!")
	l.WithContext(ctx).Info("info message!")
	l.WithContext(ctx).Warn("warn message!")
	l.WithContext(ctx).WithError(err).Error("error message!")

	immediateWrites := buffer.String()
	assert.NotContains(t, immediateWrites, "debug message!")
	assert.NotContains(t, immediateWrites, "\"hello\":\"world\"")

	assert.NotContains(t, immediateWrites, "info message!")

	assert.Contains(t, immediateWrites, "warn message!")

	assert.Contains(t, immediateWrites, "error message!")
	assert.Contains(t, immediateWrites, "\"error\":\"bad error\"")

	l.EndTracking(ctx, err)
	unwoundWrites := buffer.String()
	assert.Contains(t, unwoundWrites, "debug message!")
	assert.Contains(t, unwoundWrites, "\"hello\":\"world\"")

	assert.Contains(t, unwoundWrites, "info message!")
}
