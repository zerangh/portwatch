package eventbus_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/eventbus"
)

func TestLoggingMiddleware_WritesEntry(t *testing.T) {
	var buf bytes.Buffer
	var called bool
	h := eventbus.LoggingMiddleware(&buf, func(e eventbus.Event) {
		called = true
	})
	h(eventbus.Event{Topic: "ports"})
	if !called {
		t.Fatal("inner handler not called")
	}
	if !strings.Contains(buf.String(), "topic=ports") {
		t.Fatalf("expected log entry, got: %s", buf.String())
	}
}

func TestLoggingMiddleware_NilWriterUsesStdout(t *testing.T) {
	// should not panic
	h := eventbus.LoggingMiddleware(nil, func(e eventbus.Event) {})
	h(eventbus.Event{Topic: "ports"})
}

func TestRecoveryMiddleware_CatchesPanic(t *testing.T) {
	var buf bytes.Buffer
	h := eventbus.RecoveryMiddleware(&buf, func(e eventbus.Event) {
		panic("boom")
	})
	// should not propagate panic
	h(eventbus.Event{Topic: "ports"})
	if !strings.Contains(buf.String(), "boom") {
		t.Fatalf("expected panic message in output, got: %s", buf.String())
	}
}

func TestRecoveryMiddleware_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	var called bool
	h := eventbus.RecoveryMiddleware(&buf, func(e eventbus.Event) {
		called = true
	})
	h(eventbus.Event{Topic: "ports"})
	if !called {
		t.Fatal("inner handler not called")
	}
}
