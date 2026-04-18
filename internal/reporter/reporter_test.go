package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/state"
)

func TestReport_TextNoChanges(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatText, &buf)
	err := r.Report(state.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No port changes") {
		t.Errorf("expected no-change message, got: %s", buf.String())
	}
}

func TestReport_TextAdded(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatText, &buf)
	err := r.Report(state.Diff{Added: []int{8080, 443}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "OPENED") {
		t.Errorf("expected OPENED in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
}

func TestReport_TextRemoved(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatText, &buf)
	err := r.Report(state.Diff{Removed: []int{22}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "CLOSED") {
		t.Errorf("expected CLOSED in output, got: %s", out)
	}
}

func TestReport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatJSON, &buf)
	err := r.Report(state.Diff{Added: []int{9000}, Removed: []int{80}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"added"`) {
		t.Errorf("expected JSON keys, got: %s", out)
	}
	if !strings.Contains(out, "9000") || !strings.Contains(out, "80") {
		t.Errorf("expected ports in JSON, got: %s", out)
	}
}

func TestNew_NilWriterUsesStdout(t *testing.T) {
	// Should not panic
	r := reporter.New(reporter.FormatText, nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
