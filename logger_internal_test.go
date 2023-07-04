package unwindlogger

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestLogger_EndTrackingShouldDropEntries(t *testing.T) {
	logger := NewLogger().WithOut(io.Discard)

	ctx := context.Background()
	ctx = logger.StartTracking(ctx)

	assert.Len(t, logger.inflightContexts, 1)

	logger.WithContext(ctx).WithField("hello", "world").Info("hello")
	logger.WithContext(ctx).Warn("world")

	logger.EndTracking(ctx)
	assert.Len(t, logger.inflightContexts, 0)
}
