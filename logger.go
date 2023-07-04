package unwindlogger

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

type contextKeyType struct{}

var contextKey = contextKeyType{}

type Logger struct {
	immediateLevel   Level
	deferredLevel    Level
	out              io.Writer
	inflightContexts map[int64][]*Entry
	randSource       rand.Source
	lock             sync.Locker
	fullDefer        bool
}

func NewLogger() *Logger {
	return &Logger{
		immediateLevel:   WARN,
		deferredLevel:    INFO,
		out:              os.Stdout,
		inflightContexts: make(map[int64][]*Entry),
		randSource:       rand.NewSource(time.Now().Unix()),
		lock:             &sync.Mutex{},
		fullDefer:        false,
	}
}

func (l *Logger) WithImmediateLevel(level Level) *Logger {
	l.immediateLevel = level
	if l.deferredLevel > l.immediateLevel {
		// Deferred must always be <= immediate
		l.deferredLevel = l.immediateLevel
	}
	return l
}

func (l *Logger) WithDeferredLevel(level Level) *Logger {
	l.deferredLevel = level
	if l.deferredLevel > l.immediateLevel {
		// Immediate must always be >= deferred
		l.immediateLevel = l.deferredLevel
	}
	return l
}

func (l *Logger) WithOut(out io.Writer) *Logger {
	l.out = out
	return l
}

func (l *Logger) WithFullDefer(fullDefer bool) *Logger {
	l.fullDefer = fullDefer
	return l
}

func (l *Logger) StartTracking(ctx context.Context) context.Context {
	trackingID := l.generateTrackingID()
	l.inflightContexts[trackingID] = make([]*Entry, 0)
	return context.WithValue(ctx, contextKey, trackingID)
}

func (l *Logger) EndTracking(ctx context.Context) {
	l.EndTrackingWithError(ctx, nil)
}

func (l *Logger) EndTrackingWithError(ctx context.Context, err error) {
	if ctx == nil {
		return
	}
	trackingID, found := ctx.Value(contextKey).(int64)
	if !found {
		// Unknown context, ignore
		return
	}
	pendingEntries, found := l.inflightContexts[trackingID]
	if !found {
		// Unknown tracking ID, weird, but ignore
		_, _ = fmt.Fprintf(os.Stderr, "received entry with unknown trackingID: %d\n", trackingID)
		return
	}

	if err != nil || l.fullDefer {
		var logLevel Level
		if err != nil {
			logLevel = l.deferredLevel
		} else if l.fullDefer {
			logLevel = l.immediateLevel
		}
		for _, e := range pendingEntries {
			if logLevel <= e.level {
				e.write()
			}
		}
	}
	// Cleanup
	delete(l.inflightContexts, trackingID)
}

func (l *Logger) WithContext(ctx context.Context) *Entry {
	_, ok := ctx.Value(contextKey).(int64)
	if ok {
		return NewEntry(l, ctx)
	}
	// Not attached to any known context, just pass it through
	return NewEntry(l, nil)
}

func (l *Logger) deferEntry(e *Entry) {
	if e.ctx == nil {
		return
	}
	trackingID, found := e.ctx.Value(contextKey).(int64)
	if !found {
		// Unknown context, ignore
		return
	}
	pendingEntries, found := l.inflightContexts[trackingID]
	if !found {
		// Unknown tracking ID, weird, but ignore
		_, _ = fmt.Fprintf(os.Stderr, "received entry with unknown trackingID: %d\n", trackingID)
		return
	}

	l.inflightContexts[trackingID] = append(pendingEntries, e)
}

func (l *Logger) generateTrackingID() int64 {
	l.lock.Lock()
	i := l.randSource.Int63()
	l.lock.Unlock()
	return i
}

func (l *Logger) shouldImmediatelyLog(level Level) bool {
	return !l.fullDefer && l.immediateLevel <= level
}
