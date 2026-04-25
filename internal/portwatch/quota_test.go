package portwatch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewQuota_InvalidMax(t *testing.T) {
	_, err := NewQuota(0, time.Second, nil)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestNewQuota_InvalidWindow(t *testing.T) {
	_, err := NewQuota(5, 0, nil)
	if err == nil {
		t.Fatal("expected error for window=0")
	}
}

func TestNewQuota_NilWriterUsesStderr(t *testing.T) {
	q, err := NewQuota(1, time.Second, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestAllow_WithinQuota(t *testing.T) {
	var buf bytes.Buffer
	q, _ := NewQuota(3, time.Minute, &buf)

	for i := 0; i < 3; i++ {
		if !q.Allow() {
			t.Fatalf("call %d should be allowed", i+1)
		}
	}
}

func TestAllow_ExceedsQuota(t *testing.T) {
	var buf bytes.Buffer
	q, _ := NewQuota(2, time.Minute, &buf)

	q.Allow()
	q.Allow()

	if q.Allow() {
		t.Fatal("third call should be denied")
	}
	if !strings.Contains(buf.String(), "limit of 2") {
		t.Errorf("expected limit message in output, got: %q", buf.String())
	}
}

func TestRemaining_DecreasesWithCalls(t *testing.T) {
	q, _ := NewQuota(5, time.Minute, nil)

	if q.Remaining() != 5 {
		t.Fatalf("expected 5 remaining, got %d", q.Remaining())
	}
	q.Allow()
	q.Allow()
	if q.Remaining() != 3 {
		t.Fatalf("expected 3 remaining, got %d", q.Remaining())
	}
}

func TestReset_RestoresQuota(t *testing.T) {
	var buf bytes.Buffer
	q, _ := NewQuota(1, time.Minute, &buf)

	q.Allow() // consume the single slot
	if q.Allow() {
		t.Fatal("should be denied before reset")
	}

	q.Reset()
	if !q.Allow() {
		t.Fatal("should be allowed after reset")
	}
}

func TestAllow_WindowExpiry(t *testing.T) {
	var buf bytes.Buffer
	q, _ := NewQuota(1, 50*time.Millisecond, &buf)

	q.Allow() // consume quota
	if q.Allow() {
		t.Fatal("should be denied within window")
	}

	time.Sleep(60 * time.Millisecond)

	if !q.Allow() {
		t.Fatal("should be allowed after window expires")
	}
}
