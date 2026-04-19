package eventbus

import (
	"io"
	"os"
	"time"

	"fmt"
)

// LoggingMiddleware wraps h and writes a one-line log entry to w for each event.
func LoggingMiddleware(w io.Writer, h Handler) Handler {
	if w == nil {
		w = os.Stdout
	}
	return func(e Event) {
		start := time.Now()
		h(e)
		fmt.Fprintf(w, "[eventbus] topic=%s elapsed=%s\n", e.Topic, time.Since(start))
	}
}

// RecoveryMiddleware wraps h and recovers from any panic, writing the error to w.
func RecoveryMiddleware(w io.Writer, h Handler) Handler {
	if w == nil {
		w = os.Stderr
	}
	return func(e Event) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(w, "[eventbus] recovered panic on topic=%s: %v\n", e.Topic, r)
			}
		}()
		h(e)
	}
}
