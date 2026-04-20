package portwatch

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewLifecycle_NilWriterUsesStderr(t *testing.T) {
	l := NewLifecycle(nil)
	if l.w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestEmit_WritesToWriter(t *testing.T) {
	var buf bytes.Buffer
	l := NewLifecycle(&buf)
	l.Emit(EventStarting)
	if buf.Len() == 0 {
		t.Fatal("expected output written")
	}
	if !bytes.Contains(buf.Bytes(), []byte("starting")) {
		t.Errorf("expected 'starting' in output, got: %s", buf.String())
	}
}

func TestRegister_HandlerCalled(t *testing.T) {
	var buf bytes.Buffer
	l := NewLifecycle(&buf)

	var got []LifecycleEvent
	l.Register(func(ev LifecycleEvent, _ time.Time) {
		got = append(got, ev)
	})

	l.Emit(EventReady)
	l.Emit(EventStopped)

	if len(got) != 2 {
		t.Fatalf("expected 2 events, got %d", len(got))
	}
	if got[0] != EventReady || got[1] != EventStopped {
		t.Errorf("unexpected events: %v", got)
	}
}

func TestRegister_NilHandlerIgnored(t *testing.T) {
	var buf bytes.Buffer
	l := NewLifecycle(&buf)
	l.Register(nil)
	if len(l.handlers) != 0 {
		t.Fatal("nil handler should not be registered")
	}
}

func TestRun_EmitsLifecycleEvents(t *testing.T) {
	var buf bytes.Buffer
	l := NewLifecycle(&buf)

	var events []LifecycleEvent
	l.Register(func(ev LifecycleEvent, _ time.Time) {
		events = append(events, ev)
	})

	err := l.Run(context.Background(), func(_ context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []LifecycleEvent{EventStarting, EventReady, EventStopping, EventStopped}
	if len(events) != len(want) {
		t.Fatalf("expected %d events, got %d: %v", len(want), len(events), events)
	}
	for i, e := range want {
		if events[i] != e {
			t.Errorf("event[%d]: want %s, got %s", i, e, events[i])
		}
	}
}

func TestRun_PropagatesError(t *testing.T) {
	var buf bytes.Buffer
	l := NewLifecycle(&buf)
	sentinel := errors.New("boom")
	err := l.Run(context.Background(), func(_ context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
