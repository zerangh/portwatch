package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/state"
)

func TestNotify_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := state.Diff{
		Opened: []string{"tcp/8080", "tcp/9090"},
		Closed: []string{"tcp/3000"},
	}

	changed := n.Notify(diff)
	if !changed {
		t.Fatal("expected Notify to return true when changes exist")
	}

	out := buf.String()
	if !strings.Contains(out, "+ OPENED  tcp/8080") {
		t.Errorf("expected opened port in output, got:\n%s", out)
	}
	if !strings.Contains(out, "- CLOSED  tcp/3000") {
		t.Errorf("expected closed port in output, got:\n%s", out)
	}
}

func TestNotify_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := state.Diff{}

	changed := n.Notify(diff)
	if changed {
		t.Fatal("expected Notify to return false when no changes")
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	// Should not panic when nil writer is passed
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}
