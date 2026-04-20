package portwatch

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// LifecycleEvent represents a named stage in the watcher lifecycle.
type LifecycleEvent string

const (
	EventStarting  LifecycleEvent = "starting"
	EventReady     LifecycleEvent = "ready"
	EventStopping  LifecycleEvent = "stopping"
	EventStopped   LifecycleEvent = "stopped"
	EventScanBegin LifecycleEvent = "scan_begin"
	EventScanEnd   LifecycleEvent = "scan_end"
)

// LifecycleHandler is called when a lifecycle event fires.
type LifecycleHandler func(event LifecycleEvent, ts time.Time)

// Lifecycle tracks and broadcasts watcher lifecycle events.
type Lifecycle struct {
	w        io.Writer
	handlers []LifecycleHandler
}

// NewLifecycle creates a Lifecycle that logs events to w.
// If w is nil, os.Stderr is used.
func NewLifecycle(w io.Writer) *Lifecycle {
	if w == nil {
		w = os.Stderr
	}
	return &Lifecycle{w: w}
}

// Register adds a handler that will be called on every lifecycle event.
func (l *Lifecycle) Register(h LifecycleHandler) {
	if h != nil {
		l.handlers = append(l.handlers, h)
	}
}

// Emit fires event, logs it, and notifies all registered handlers.
func (l *Lifecycle) Emit(event LifecycleEvent) {
	ts := time.Now()
	fmt.Fprintf(l.w, "[lifecycle] %s at %s\n", event, ts.Format(time.RFC3339))
	for _, h := range l.handlers {
		h(event, ts)
	}
}

// Run wraps fn with starting/ready/stopping/stopped events and honours ctx.
func (l *Lifecycle) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	l.Emit(EventStarting)
	l.Emit(EventReady)
	defer func() {
		l.Emit(EventStopping)
		l.Emit(EventStopped)
	}()
	return fn(ctx)
}
