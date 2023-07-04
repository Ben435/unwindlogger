package unwindlogger

import (
	"context"
	"fmt"
	"io"
	"testing"
)

func BenchmarkLogger_ImmediateLogging(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		logger := NewLogger().
			WithImmediateLevel(DEBUG).
			WithOut(io.Discard)
		ctx := context.Background()
		for pb.Next() {
			ctx := logger.StartTracking(ctx)
			logger.WithContext(context.Background()).Debug("aaa")
			logger.EndTracking(ctx)
		}
	})
}

func BenchmarkLogger_DeferredLogging(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		logger := NewLogger().
			WithImmediateLevel(DEBUG).
			WithFullDefer(true).
			WithOut(io.Discard)
		ctx := context.Background()
		for pb.Next() {
			ctx := logger.StartTracking(ctx)
			logger.WithContext(ctx).Debug("aaa")
			logger.EndTracking(ctx)
		}
	})
}

func BenchmarkLogger_MixedLoggingWithError(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		logger := NewLogger().
			WithDeferredLevel(INFO).
			WithImmediateLevel(WARN).
			WithFullDefer(false).
			WithOut(io.Discard)
		err := fmt.Errorf("err")
		ctx := context.Background()
		for pb.Next() {
			ctx := logger.StartTracking(ctx)
			logger.WithContext(ctx).Debug("aaa")
			logger.WithContext(ctx).Info("bbb")
			logger.WithContext(ctx).Warn("ccc")
			logger.EndTrackingWithError(ctx, err)
		}
	})
}

func BenchmarkLogger_MixedLoggingWithoutError(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		logger := NewLogger().
			WithDeferredLevel(INFO).
			WithImmediateLevel(WARN).
			WithFullDefer(false).
			WithOut(io.Discard)
		ctx := context.Background()
		for pb.Next() {
			ctx := logger.StartTracking(ctx)
			logger.WithContext(ctx).Debug("aaa")
			logger.WithContext(ctx).Info("bbb")
			logger.WithContext(ctx).Warn("ccc")
			logger.EndTracking(ctx)
		}
	})
}

func logStuff(logger *Logger, ctx context.Context) {
	logger.WithContext(ctx).WithField("hello", "world").Debug("aaa")
	logger.WithContext(ctx).Info("bbb")
	logger.WithContext(ctx).Warn("ccc")

	logger.WithContext(ctx).WithField("hello", "world").Debug("aaa")
	logger.WithContext(ctx).Info("bbb")
	logger.WithContext(ctx).Warn("ccc")
	logger.WithContext(ctx).WithError(fmt.Errorf("an error")).Error("ddd")

	logger.WithContext(ctx).WithField("hello", "world").Debug("aaa")
	logger.WithContext(ctx).Info("bbb")
	logger.WithContext(ctx).Warn("ccc")
}

func BenchmarkLogger_MixedLoggingWithError10Logs(b *testing.B) {
	logger := NewLogger().
		WithDeferredLevel(INFO).
		WithImmediateLevel(WARN).
		WithFullDefer(false).
		WithOut(io.Discard)
	err := fmt.Errorf("err")
	ctx := context.Background()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx1 := logger.StartTracking(ctx)
			logStuff(logger, ctx1)

			ctx2 := logger.StartTracking(ctx)
			logStuff(logger, ctx2)

			ctx3 := logger.StartTracking(ctx)
			logStuff(logger, ctx3)

			logger.EndTrackingWithError(ctx1, err)
			logger.EndTrackingWithError(ctx2, err)
			logger.EndTrackingWithError(ctx3, err)
		}
	})
}
