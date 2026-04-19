package signal_test

import (
	"bytes"
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/user/portwatch/internal/signal"
)

func TestNew_DefaultsToStderr(t *testing.T) {
	h := signal.New(nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestNotify_CancelledByParent(t *testing.T) {
	buf := &bytes.Buffer{}
	h := signal.New(buf)
	parent, parentCancel := context.WithCancel(context.Background())

	ctx, stop := h.Notify(parent)
	defer stop()

	parentCancel()

	select {
	case <-ctx.Done():
		// ok
	case <-time.After(time.Second):
		t.Fatal("context not cancelled after parent cancel")
	}
}

func TestNotify_CancelledBySignal(t *testing.T) {
	buf := &bytes.Buffer{}
	h := signal.New(buf, syscall.SIGUSR1)

	ctx, stop := h.Notify(context.Background())
	defer stop()

	// send the signal to ourselves
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)

	select {
	case <-ctx.Done():
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("context not cancelled after signal")
	}

	if buf.Len() == 0 {
		t.Error("expected shutdown message written")
	}
}

func TestWait_ReturnsWhenDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan struct{})
	go func() {
		signal.Wait(ctx)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Wait did not return")
	}
}
