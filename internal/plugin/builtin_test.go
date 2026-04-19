package plugin_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/plugin"
)

func TestLogHandler_WritesChanges(t *testing.T) {
	var buf bytes.Buffer
	h := plugin.LogHandler(&buf)
	err := h(plugin.Event{Host: "localhost", Added: []int{22, 80}, Removed: []int{443}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "host=localhost") {
		t.Errorf("expected host in output, got: %s", out)
	}
	if !strings.Contains(out, "added=") {
		t.Errorf("expected added in output, got: %s", out)
	}
	if !strings.Contains(out, "removed=") {
		t.Errorf("expected removed in output, got: %s", out)
	}
}

func TestLogHandler_NoChanges_Silent(t *testing.T) {
	var buf bytes.Buffer
	h := plugin.LogHandler(&buf)
	_ = h(plugin.Event{Host: "localhost"})
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty event, got: %s", buf.String())
	}
}

func TestThresholdHandler_BelowThreshold(t *testing.T) {
	called := false
	action := func(e plugin.Event) error { called = true; return nil }
	h := plugin.ThresholdHandler(5, action)
	_ = h(plugin.Event{Added: []int{80, 443}})
	if called {
		t.Error("action should not be called below threshold")
	}
}

func TestThresholdHandler_MeetsThreshold(t *testing.T) {
	called := false
	action := func(e plugin.Event) error { called = true; return nil }
	h := plugin.ThresholdHandler(2, action)
	_ = h(plugin.Event{Added: []int{80, 443}})
	if !called {
		t.Error("action should be called at threshold")
	}
}
