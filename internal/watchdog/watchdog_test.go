package watchdog

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNew_InvalidMaxAge(t *testing.T) {
	_, err := New(0, nil)
	if err == nil {
		t.Fatal("expected error for zero maxAge")
	}
}

func TestNew_NilWriterUsesStderr(t *testing.T) {
	wd, err := New(time.Minute, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wd.writer == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestIsStale_NeverPinged(t *testing.T) {
	wd, _ := New(time.Minute, nil)
	if !wd.IsStale() {
		t.Fatal("expected stale when never pinged")
	}
}

func TestIsStale_AfterPing(t *testing.T) {
	wd, _ := New(time.Minute, nil)
	wd.Ping()
	if wd.IsStale() {
		t.Fatal("expected not stale immediately after ping")
	}
}

func TestIsStale_Expired(t *testing.T) {
	wd, _ := New(10*time.Millisecond, nil)
	wd.Ping()
	time.Sleep(20 * time.Millisecond)
	if !wd.IsStale() {
		t.Fatal("expected stale after maxAge elapsed")
	}
}

func TestAge_NeverPinged(t *testing.T) {
	wd, _ := New(time.Minute, nil)
	if wd.Age() != -1 {
		t.Fatal("expected -1 when never pinged")
	}
}

func TestAge_AfterPing(t *testing.T) {
	wd, _ := New(time.Minute, nil)
	wd.Ping()
	if wd.Age() < 0 {
		t.Fatal("expected non-negative age after ping")
	}
}

func TestCheck_WritesWarningWhenStale(t *testing.T) {
	var buf bytes.Buffer
	wd, _ := New(10*time.Millisecond, &buf)
	time.Sleep(20 * time.Millisecond)
	wd.Check()
	if !strings.Contains(buf.String(), "WARNING") {
		t.Fatalf("expected WARNING in output, got: %q", buf.String())
	}
}

func TestCheck_SilentWhenHealthy(t *testing.T) {
	var buf bytes.Buffer
	wd, _ := New(time.Minute, &buf)
	wd.Ping()
	wd.Check()
	if buf.Len() != 0 {
		t.Fatalf("expected no output when healthy, got: %q", buf.String())
	}
}
