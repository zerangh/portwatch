package portwatch

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/example/portwatch/internal/portdiff"
)

func TestNewHooks_NilWriterUsesStderr(t *testing.T) {
	h := NewHooks(nil)
	if h.w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestRegister_IncreasesLen(t *testing.T) {
	h := NewHooks(nil)
	if h{
		t.Fatalf("expected 0 handlers, got %d", h.Len())
	}
	h.Register(func(_ HookEvent, _ *portdiff.Diff, _ error) {})
	h.Register(func(_ HookEvent, _ *portdiff.Diff, _ error) {})
	if h.Len() != 2 {
		t.Fatalf("expected 2 handlers, got %d", h.Len())
	}
}

func TestRegister_NilFuncIgnored(t *testing.T) {
	h := NewHooks(nil)
	h.Register(nil)
	if h.Len() != 0 {
		t.Fatalf("expected 0 handlers after nil register, got %d", h.Len())
	}
}

func TestFire_InvokesAllHandlers(t *testing.T) {
	h := NewHooks(nil)
	var called []HookEvent
	for i := 0; i < 3; i++ {
		h.Register(func(ev HookEvent, _ *portdiff.Diff, _ error) {
			called = append(called, ev)
		})
	}
	h.Fire(HookOnChange, nil, nil)
	if len(called) != 3 {
		t.Fatalf("expected 3 invocations, got %d", len(called))
	}
	for _, ev := range called {
		if ev != HookOnChange {
			t.Errorf("unexpected event %q", ev)
		}
	}
}

func TestFire_PassesDiffAndError(t *testing.T) {
	h := NewHooks(nil)
	sentinelErr := errors.New("boom")
	var gotErr error
	h.Register(func(_ HookEvent, _ *portdiff.Diff, err error) {
		gotErr = err
	})
	h.Fire(HookOnError, nil, sentinelErr)
	if gotErr != sentinelErr {
		t.Errorf("expected sentinel error, got %v", gotErr)
	}
}

func TestFire_RecoversPanic(t *testing.T) {
	var buf bytes.Buffer
	h := NewHooks(&buf)
	h.Register(func(_ HookEvent, _ *portdiff.Diff, _ error) {
		panic("test panic")
	})
	// Should not panic the caller.
	h.Fire(HookAfterScan, nil, nil)
	if !strings.Contains(buf.String(), "test panic") {
		t.Errorf("expected panic message in output, got: %q", buf.String())
	}
}

func TestEventNames_ContainsAllEvents(t *testing.T) {
	names := EventNames()
	for _, ev := range []HookEvent{HookBeforeScan, HookAfterScan, HookOnChange, HookOnError} {
		if !strings.Contains(names, string(ev)) {
			t.Errorf("EventNames missing %q", ev)
		}
	}
}
