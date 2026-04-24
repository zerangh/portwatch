package portwatch

import (
	"fmt"
	"io"
	"os"
	"time"
)

// ScanTimingMiddleware returns a LifecycleHandler that records durations
// between scan_begin and scan_end events, writing them to w.
func ScanTimingMiddleware(w io.Writer) LifecycleHandler {
	if w == nil {
		w = os.Stdout
	}
	var start time.Time
	return func(ev LifecycleEvent, ts time.Time) {
		switch ev {
		case EventScanBegin:
			start = ts
		case EventScanEnd:
			if !start.IsZero() {
				fmt.Fprintf(w, "[timing] scan completed in %s\n", ts.Sub(start).Round(time.Millisecond))
				start = time.Time{}
			}
		}
	}
}

// UptimeMiddleware returns a LifecycleHandler that prints elapsed time
// since EventReady when EventStopping fires.
func UptimeMiddleware(w io.Writer) LifecycleHandler {
	if w == nil {
		w = os.Stdout
	}
	var readyAt time.Time
	return func(ev LifecycleEvent, ts time.Time) {
		switch ev {
		case EventReady:
			readyAt = ts
		case EventStopping:
			if !readyAt.IsZero() {
				fmt.Fprintf(w, "[uptime] ran for %s\n", ts.Sub(readyAt).Round(time.Second))
			}
		}
	}
}

// ChainMiddleware combines multiple LifecycleHandlers into one, invoking
// each handler in order for every event received.
func ChainMiddleware(handlers ...LifecycleHandler) LifecycleHandler {
	return func(ev LifecycleEvent, ts time.Time) {
		for _, h := range handlers {
			if h != nil {
				h(ev, ts)
			}
		}
	}
}
