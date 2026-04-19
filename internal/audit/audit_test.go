package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNew_NilWriterUsesStdout(t *testing.T) {
	l := New(nil, false)
	if l.w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestLog_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, false)
	if err := l.Log(LevelInfo, "scan.start", "ports 1-1024"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO in output, got: %s", out)
	}
	if !strings.Contains(out, "scan.start") {
		t.Errorf("expected event in output, got: %s", out)
	}
}

func TestLog_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, true)
	if err := l.Info("scan.complete", "3 ports open"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var e Entry
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if e.Level != LevelInfo {
		t.Errorf("expected INFO, got %s", e.Level)
	}
	if e.Event != "scan.complete" {
		t.Errorf("expected scan.complete, got %s", e.Event)
	}
}

func TestWarn_ContainsWARN(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, false)
	_ = l.Warn("port.changed", "added: 8080")
	if !strings.Contains(buf.String(), "WARN") {
		t.Errorf("expected WARN in output")
	}
}

func TestError_ContainsERROR(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, false)
	_ = l.Error("scan.failed", "connection refused")
	if !strings.Contains(buf.String(), "ERROR") {
		t.Errorf("expected ERROR in output")
	}
}
