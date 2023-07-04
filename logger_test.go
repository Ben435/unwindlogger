package unwindlogger_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Ben435/unwindlogger"
)

func TestLogger_DirectLogging(t *testing.T) {
	buffer := &bytes.Buffer{}
	ctx := context.Background()

	l := unwindlogger.NewLogger().
		WithImmediateLevel(unwindlogger.DEBUG).
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
		WithImmediateLevel(unwindlogger.WARN).
		WithDeferredLevel(unwindlogger.INFO).
		WithOut(buffer)

	ctx = l.StartTracking(ctx)

	l.WithContext(ctx).Debug("debug message!")
	l.WithContext(ctx).WithField("hello", "world").Info("info message!")
	l.WithContext(ctx).Warn("warn message!")
	l.WithContext(ctx).WithError(err).Error("error message!")

	immediateWrites := buffer.String()
	assert.NotContains(t, immediateWrites, "debug message!")

	assert.NotContains(t, immediateWrites, "info message!")
	assert.NotContains(t, immediateWrites, "\"hello\":\"world\"")

	assert.Contains(t, immediateWrites, "warn message!")

	assert.Contains(t, immediateWrites, "error message!")
	assert.Contains(t, immediateWrites, "\"error\":\"bad error\"")

	l.EndTrackingWithError(ctx, err)
	unwoundWrites := buffer.String()
	assert.NotContains(t, immediateWrites, "debug message!")

	assert.Contains(t, unwoundWrites, "info message!")
	assert.Contains(t, unwoundWrites, "\"hello\":\"world\"")
}

func TestLogger_UnwindLoggingWithFullDeferAndWhenEndingTrackingWithError(t *testing.T) {
	buffer := &bytes.Buffer{}
	ctx := context.Background()
	err := fmt.Errorf("bad error")

	l := unwindlogger.NewLogger().
		WithImmediateLevel(unwindlogger.WARN).
		WithDeferredLevel(unwindlogger.INFO).
		WithFullDefer(true).
		WithOut(buffer)

	ctx = l.StartTracking(ctx)

	l.WithContext(ctx).Debug("debug message!")
	l.WithContext(ctx).WithField("hello", "world").Info("info message!")
	l.WithContext(ctx).Warn("warn message!")
	l.WithContext(ctx).WithError(err).Error("error message!")

	immediateWrites := buffer.String()
	assert.Empty(t, immediateWrites)

	l.EndTrackingWithError(ctx, err)
	unwoundWrites := buffer.String()
	assert.NotContains(t, unwoundWrites, "debug message!")

	assert.Contains(t, unwoundWrites, "info message!")
	assert.Contains(t, unwoundWrites, "\"hello\":\"world\"")

	assert.Contains(t, unwoundWrites, "warn message!")

	assert.Contains(t, unwoundWrites, "error message!")
	assert.Contains(t, unwoundWrites, "\"error\":\"bad error\"")
}

func TestLogger_UnwindLoggingWithFullDeferAndWhenEndingTrackingWithoutError(t *testing.T) {
	buffer := &bytes.Buffer{}
	ctx := context.Background()

	l := unwindlogger.NewLogger().
		WithImmediateLevel(unwindlogger.WARN).
		WithDeferredLevel(unwindlogger.INFO).
		WithFullDefer(true).
		WithOut(buffer)

	ctx = l.StartTracking(ctx)

	l.WithContext(ctx).Debug("debug message!")
	l.WithContext(ctx).WithField("hello", "world").Info("info message!")
	l.WithContext(ctx).Warn("warn message!")
	l.WithContext(ctx).Error("error message!")

	immediateWrites := buffer.String()
	assert.Empty(t, immediateWrites)

	l.EndTracking(ctx)
	unwoundWrites := buffer.String()
	assert.NotContains(t, unwoundWrites, "debug message!")

	assert.NotContains(t, unwoundWrites, "info message!")
	assert.NotContains(t, unwoundWrites, "\"hello\":\"world\"")

	assert.Contains(t, unwoundWrites, "warn message!")

	assert.Contains(t, unwoundWrites, "error message!")
}
