package portwatch

import (
	"bytes"
	"testing"
	"time"
)

func TestNewSampler_Defaults(t *testing.T) {
	s := NewSampler(0, nil)
	if s.maxSize != 100 {
		t.Fatalf("expected maxSize 100, got %d", s.maxSize)
	}
	if s.w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestSampler_Mean_Empty(t *testing.T) {
	s := NewSampler(10, &bytes.Buffer{})
	if got := s.Mean(); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestSampler_Mean_Computed(t *testing.T) {
	s := NewSampler(10, &bytes.Buffer{})
	s.Record(10 * time.Millisecond)
	s.Record(20 * time.Millisecond)
	s.Record(30 * time.Millisecond)
	want := 20 * time.Millisecond
	if got := s.Mean(); got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestSampler_P95_Empty(t *testing.T) {
	s := NewSampler(10, &bytes.Buffer{})
	if got := s.P95(); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestSampler_P95_SingleSample(t *testing.T) {
	s := NewSampler(10, &bytes.Buffer{})
	s.Record(42 * time.Millisecond)
	if got := s.P95(); got != 42*time.Millisecond {
		t.Fatalf("expected 42ms, got %v", got)
	}
}

func TestSampler_P95_LargeWindow(t *testing.T) {
	s := NewSampler(100, &bytes.Buffer{})
	for i := 1; i <= 100; i++ {
		s.Record(time.Duration(i) * time.Millisecond)
	}
	// p95 index = int(100*0.95)-1 = 94 → 95ms (0-based sorted)
	want := 95 * time.Millisecond
	if got := s.P95(); got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestSampler_EvictsOldest(t *testing.T) {
	s := NewSampler(3, &bytes.Buffer{})
	s.Record(1 * time.Millisecond)
	s.Record(2 * time.Millisecond)
	s.Record(3 * time.Millisecond)
	s.Record(4 * time.Millisecond) // evicts 1ms
	if got := s.Len(); got != 3 {
		t.Fatalf("expected len 3, got %d", got)
	}
	// mean should be (2+3+4)/3 = 3ms
	want := 3 * time.Millisecond
	if got := s.Mean(); got != want {
		t.Fatalf("expected mean %v, got %v", want, got)
	}
}

func TestSampler_Len(t *testing.T) {
	s := NewSampler(10, &bytes.Buffer{})
	for i := 0; i < 5; i++ {
		s.Record(time.Millisecond)
	}
	if got := s.Len(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}
