package portdiff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portdiff"
)

func TestPrint_TextNoChanges(t *testing.T) {
	var buf bytes.Buffer
	err := portdiff.Print(&buf, portdiff.Result{}, portdiff.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no port changes") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestPrint_TextAdded(t *testing.T) {
	var buf bytes.Buffer
	r := portdiff.Result{Added: []int{80, 443}}
	if err := portdiff.Print(&buf, r, portdiff.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "added") {
		t.Fatalf("expected 'added' in output, got %q", out)
	}
	if !strings.Contains(out, "80") || !strings.Contains(out, "443") {
		t.Fatalf("expected ports in output, got %q", out)
	}
}

func TestPrint_TextRemoved(t *testing.T) {
	var buf bytes.Buffer
	r := portdiff.Result{Removed: []int{22}}
	if err := portdiff.Print(&buf, r, portdiff.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "removed") {
		t.Fatalf("expected 'removed' in output, got %q", buf.String())
	}
}

func TestPrint_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	r := portdiff.Result{Added: []int{8080}, Removed: []int{22}}
	if err := portdiff.Print(&buf, r, portdiff.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "added") || !strings.Contains(out, "removed") {
		t.Fatalf("expected JSON keys in output, got %q", out)
	}
}

func TestPrint_NilWriterUsesStdout(t *testing.T) {
	// Should not panic when writer is nil.
	err := portdiff.Print(nil, portdiff.Result{}, portdiff.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPrint_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := portdiff.Print(&buf, portdiff.Result{}, portdiff.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}
