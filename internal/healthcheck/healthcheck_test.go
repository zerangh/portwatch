package healthcheck

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNew_NilWriterUsesStdout(t *testing.T) {
	c := New(nil)
	if c.w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestRecordScan_SetsHealthy(t *testing.T) {
	c := New(nil)
	now := time.Now()
	c.RecordScan(now)
	s := c.Status()
	if !s.Healthy {
		t.Error("expected healthy after scan")
	}
	if s.ScanCount != 1 {
		t.Errorf("expected scan count 1, got %d", s.ScanCount)
	}
	if !s.LastScan.Equal(now) {
		t.Error("expected last scan time to match")
	}
}

func TestRecordError_SetsUnhealthy(t *testing.T) {
	c := New(nil)
	c.RecordError(errors.New("timeout"))
	s := c.Status()
	if s.Healthy {
		t.Error("expected unhealthy after error")
	}
	if s.ErrorCount != 1 {
		t.Errorf("expected error count 1, got %d", s.ErrorCount)
	}
	if s.LastError != "timeout" {
		t.Errorf("unexpected last error: %s", s.LastError)
	}
}

func TestPrint_HealthyOutput(t *testing.T) {
	var buf bytes.Buffer
	c := New(&buf)
	c.RecordScan(time.Now())
	c.Print()
	if !strings.Contains(buf.String(), "OK") {
		t.Errorf("expected OK in output: %s", buf.String())
	}
}

func TestPrint_DegradedOutput(t *testing.T) {
	var buf bytes.Buffer
	c := New(&buf)
	c.RecordError(errors.New("conn refused"))
	c.Print()
	out := buf.String()
	if !strings.Contains(out, "DEGRADED") {
		t.Errorf("expected DEGRADED in output: %s", out)
	}
	if !strings.Contains(out, "conn refused") {
		t.Errorf("expected error in output: %s", out)
	}
}

func TestPrint_NeverScanned(t *testing.T) {
	var buf bytes.Buffer
	c := New(&buf)
	c.Print()
	if !strings.Contains(buf.String(), "never") {
		t.Errorf("expected 'never' when no scan recorded: %s", buf.String())
	}
}
