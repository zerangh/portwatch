package plugin_test

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/user/portwatch/internal/plugin"
)

func TestRegisterAndDispatch(t *testing.T) {
	r := plugin.New()
	var called int32
	r.Register("counter", func(e plugin.Event) error {
		atomic.AddInt32(&called, 1)
		return nil
	})
	errs := r.Dispatch(plugin.Event{Host: "localhost", Added: []int{8080}})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if atomic.LoadInt32(&called) != 1 {
		t.Fatalf("expected handler called once, got %d", called)
	}
}

func TestUnregister(t *testing.T) {
	r := plugin.New()
	r.Register("h", func(e plugin.Event) error { return nil })
	if r.Len() != 1 {
		t.Fatalf("expected 1 handler")
	}
	r.Unregister("h")
	if r.Len() != 0 {
		t.Fatalf("expected 0 handlers after unregister")
	}
}

func TestDispatch_CollectsErrors(t *testing.T) {
	r := plugin.New()
	r.Register("bad", func(e plugin.Event) error { return errors.New("boom") })
	r.Register("ok", func(e plugin.Event) error { return nil })
	errs := r.Dispatch(plugin.Event{})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestDispatch_Empty(t *testing.T) {
	r := plugin.New()
	errs := r.Dispatch(plugin.Event{Host: "localhost"})
	if len(errs) != 0 {
		t.Fatalf("expected no errors on empty registry")
	}
}

func TestRegister_Overwrite(t *testing.T) {
	r := plugin.New()
	r.Register("h", func(e plugin.Event) error { return errors.New("first") })
	r.Register("h", func(e plugin.Event) error { return nil })
	if r.Len() != 1 {
		t.Fatalf("expected 1 handler after overwrite")
	}
	errs := r.Dispatch(plugin.Event{})
	if len(errs) != 0 {
		t.Fatalf("expected overwritten handler to return no error")
	}
}
