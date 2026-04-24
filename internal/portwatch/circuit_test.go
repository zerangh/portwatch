package portwatch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewCircuit_InvalidArgs(t *testing.T) {
	_, err := NewCircuit(0, time.Second, nil)
	if err == nil {
		t.Fatal("expected error for maxFailures=0")
	}
	_, err = NewCircuit(1, 0, nil)
	if err == nil {
		t.Fatal("expected error for resetAfter=0")
	}
}

func TestNewCircuit_NilWriterUsesStderr(t *testing.T) {
	c, err := NewCircuit(2, time.Second, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestCircuit_InitiallyClosed(t *testing.T) {
	c, _ := NewCircuit(3, time.Second, &bytes.Buffer{})
	if !c.Allow() {
		t.Fatal("expected circuit to allow in closed state")
	}
	if c.State() != StateClosed {
		t.Fatalf("expected closed, got %s", c.State())
	}
}

func TestCircuit_OpensAfterMaxFailures(t *testing.T) {
	var buf bytes.Buffer
	c, _ := NewCircuit(3, time.Second, &buf)
	for i := 0; i < 3; i++ {
		c.RecordFailure()
	}
	if c.State() != StateOpen {
		t.Fatalf("expected open, got %s", c.State())
	}
	if c.Allow() {
		t.Fatal("expected circuit to block in open state")
	}
	if !strings.Contains(buf.String(), "opened") {
		t.Error("expected opened message in output")
	}
}

func TestCircuit_HalfOpenAfterReset(t *testing.T) {
	var buf bytes.Buffer
	c, _ := NewCircuit(1, 10*time.Millisecond, &buf)
	c.RecordFailure()
	if c.State() != StateOpen {
		t.Fatalf("expected open, got %s", c.State())
	}
	time.Sleep(20 * time.Millisecond)
	if !c.Allow() {
		t.Fatal("expected circuit to allow in half-open state")
	}
	if c.State() != StateHalfOpen {
		t.Fatalf("expected half-open, got %s", c.State())
	}
}

func TestCircuit_RecoveryCloses(t *testing.T) {
	var buf bytes.Buffer
	c, _ := NewCircuit(1, 10*time.Millisecond, &buf)
	c.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	c.Allow() // transition to half-open
	c.RecordSuccess()
	if c.State() != StateClosed {
		t.Fatalf("expected closed after success, got %s", c.State())
	}
	if !strings.Contains(buf.String(), "closed after recovery") {
		t.Error("expected recovery message in output")
	}
}

func TestCircuit_SuccessResetFailureCount(t *testing.T) {
	c, _ := NewCircuit(3, time.Second, &bytes.Buffer{})
	c.RecordFailure()
	c.RecordFailure()
	c.RecordSuccess()
	if c.State() != StateClosed {
		t.Fatalf("expected closed, got %s", c.State())
	}
	// Two more failures should not open (count reset)
	c.RecordFailure()
	c.RecordFailure()
	if c.State() != StateClosed {
		t.Fatalf("expected still closed, got %s", c.State())
	}
}
