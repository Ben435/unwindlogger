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

	assert.Len(t, logger.inflightTrackingIDs, initialTrackingEntriesCapacity)
	assert.Len(t, logger.inflightTrackingIDs[0], 0)

	logger.WithContext(ctx).WithField("hello", "world").Info("hello")
	logger.WithContext(ctx).Warn("world")

	assert.Len(t, logger.inflightTrackingIDs[0], 1)

	logger.EndTracking(ctx)
	assert.Len(t, logger.inflightTrackingIDs, initialTrackingEntriesCapacity) // Retain capacity
	assert.Len(t, logger.inflightTrackingIDs[0], 0)                           // But drop the inflight items
}
