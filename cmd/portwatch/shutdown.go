package main

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/signal"
)

const shutdownTimeout = 5 * time.Second

// shutdownContext returns a context that is cancelled on SIGINT/SIGTERM.
// The returned stop function must be deferred by the caller.
func shutdownContext(w io.Writer) (context.Context, context.CancelFunc) {
	if w == nil {
		w = os.Stderr
	}
	h := signal.New(w)
	return h.Notify(context.Background())
}

// waitForShutdown blocks until ctx is done, then waits up to
// shutdownTimeout for the provided cleanup function to complete.
func waitForShutdown(ctx context.Context, cleanup func()) {
	<-ctx.Done()
	done := make(chan struct{})
	go func() {
		cleanup()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(shutdownTimeout):
		_, _ = io.WriteString(os.Stderr, "portwatch: shutdown timed out\n")
	}
}
