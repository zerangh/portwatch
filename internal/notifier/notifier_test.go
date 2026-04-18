package notifier_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/state"
)

func TestStdoutNotifier_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.NewStdout(&buf)
	if err := n.Notify(nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestStdoutNotifier_AddedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.NewStdout(&buf)
	added := []state.Port{{Number: 8080, Proto: "tcp"}}
	if err := n.Notify(added, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "NEW ports open") {
		t.Errorf("expected 'NEW ports open' in output, got %q", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got %q", out)
	}
}

func TestStdoutNotifier_RemovedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.NewStdout(&buf)
	removed := []state.Port{{Number: 22, Proto: "tcp"}}
	if err := n.Notify(nil, removed); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Ports closed") {
		t.Errorf("expected 'Ports closed' in output, got %q", out)
	}
}

func TestNewStdout_NilUsesStdout(t *testing.T) {
	n := notifier.NewStdout(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
