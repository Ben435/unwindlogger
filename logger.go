package unwindlogger

import (
	"context"
	"io"
	"os"
	"sync"
)

type contextKeyType struct{}

var contextKey = contextKeyType{}

type Logger struct {
	immediateLevel      Level
	deferredLevel       Level
	out                 io.Writer
	inflightTrackingIDs [][]*Entry
	fullDefer           bool
	trackingIDPool      *sync.Pool
	lock                sync.Locker
}

const initialIDs = 10
const initialEntries = 10

func NewLogger() *Logger {
	inflightTrackingIDs := make([][]*Entry, initialIDs)
	trackingIDPool := &sync.Pool{}

	for i := 0; i < initialIDs; i++ {
		inflightTrackingIDs[i] = make([]*Entry, 0, initialEntries)
		trackingIDPool.Put(i)
	}

	return &Logger{
		immediateLevel:      WARN,
		deferredLevel:       INFO,
		out:                 os.Stdout,
		inflightTrackingIDs: inflightTrackingIDs,
		trackingIDPool:      trackingIDPool,
		fullDefer:           false,
		lock:                &sync.Mutex{},
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
	l.lock.Lock()
	defer l.lock.Unlock()
	trackingID, ok := l.trackingIDPool.Get().(int)
	if !ok {
		trackingID = len(l.inflightTrackingIDs)
		l.inflightTrackingIDs = append(l.inflightTrackingIDs, []*Entry{})
	}

	return context.WithValue(ctx, contextKey, trackingID)
}

func (l *Logger) EndTracking(ctx context.Context) {
	l.EndTrackingWithError(ctx, nil)
}

func (l *Logger) EndTrackingWithError(ctx context.Context, err error) {
	if ctx == nil {
		return
	}
	trackingID, found := ctx.Value(contextKey).(int)
	if !found {
		// Unknown context, ignore
		return
	}

	if err != nil || l.fullDefer {
		pendingEntries := l.inflightTrackingIDs[trackingID]

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
	l.inflightTrackingIDs[trackingID] = l.inflightTrackingIDs[trackingID][:0]
	l.trackingIDPool.Put(trackingID)
}

func (l *Logger) WithContext(ctx context.Context) *Entry {
	return NewEntry(l, ctx)
}

func (l *Logger) deferEntry(e *Entry) {
	if e.ctx == nil || l.deferredLevel > e.level {
		return
	}
	trackingID, found := e.ctx.Value(contextKey).(int)
	if !found {
		// Unknown context, ignore
		return
	}

	l.inflightTrackingIDs[trackingID] = append(l.inflightTrackingIDs[trackingID], e)
}

func (l *Logger) shouldImmediatelyLog(level Level) bool {
	return !l.fullDefer && l.immediateLevel <= level
}
