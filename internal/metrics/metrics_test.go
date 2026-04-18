package metrics_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestNew_ZeroValues(t *testing.T) {
	c := metrics.New()
	s := c.Snapshot()
	if s.Scans != 0 || s.Changes != 0 || s.Errors != 0 {
		t.Fatalf("expected zero counters, got %+v", s)
	}
}

func TestRecord_NoChange(t *testing.T) {
	c := metrics.New()
	c.Record(false, nil)
	s := c.Snapshot()
	if s.Scans != 1 || s.Changes != 0 || s.Errors != 0 {
		t.Fatalf("unexpected counters: %+v", s)
	}
	if s.LastScan.IsZero() {
		t.Fatal("LastScan should be set")
	}
}

func TestRecord_WithChange(t *testing.T) {
	c := metrics.New()
	c.Record(true, nil)
	s := c.Snapshot()
	if s.Changes != 1 {
		t.Fatalf("expected 1 change, got %d", s.Changes)
	}
	if s.LastChange.IsZero() {
		t.Fatal("LastChange should be set")
	}
}

func TestRecord_WithError(t *testing.T) {
	c := metrics.New()
	c.Record(false, errors.New("scan failed"))
	s := c.Snapshot()
	if s.Errors != 1 || s.Changes != 0 {
		t.Fatalf("unexpected counters: %+v", s)
	}
}

func TestRecord_Concurrent(t *testing.T) {
	c := metrics.New()
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			c.Record(true, nil)
			done <- struct{}{}
		}()
	}
	timeout := time.After(2 * time.Second)
	for i := 0; i < 50; i++ {
		select {
		case <-done:
		case <-timeout:
			t.Fatal("timed out waiting for goroutines")
		}
	}
	if s := c.Snapshot(); s.Scans != 50 {
		t.Fatalf("expected 50 scans, got %d", s.Scans)
	}
}
