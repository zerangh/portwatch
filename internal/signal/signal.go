// Package signal provides OS signal handling for graceful shutdown.
package signal

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// Handler listens for OS signals and cancels a context.
type Handler struct {
	writer io.Writer
	sigs   []os.Signal
}

// New creates a Handler that reacts to the given signals.
// If sigs is empty, SIGINT and SIGTERM are used.
func New(w io.Writer, sigs ...os.Signal) *Handler {
	if w == nil {
		w = os.Stderr
	}
	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	return &Handler{writer: w, sigs: sigs}
}

// Notify returns a context that is cancelled when one of the registered
// signals is received. The caller must invoke the returned stop function
// to release resources when the context is no longer needed.
func (h *Handler) Notify(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, h.sigs...)

	go func() {
		defer signal.Stop(ch)
		select {
		case sig := <-ch:
			_, = io.WriteString(h.writer, "portwatch: received signal "+sig.String()+", shutting down\n")
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}

// Wait blocks until the context is done.
func Wait(ctx context.Context) {
	<-ctx.Done()
}
