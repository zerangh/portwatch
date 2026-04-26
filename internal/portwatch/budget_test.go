package portwatch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewBudget_InvalidWindow(t *testing.T) {
	_, err := NewBudget(0, time.Second, nil)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewBudget_InvalidThreshold(t *testing.T) {
	_, err := NewBudget(time.Minute, 0, nil)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNewBudget_NilWriterUsesStderr(t *testing.T) {
	b, err := NewBudget(time.Minute, time.Second, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestRecord_WithinBudget(t *testing.T) {
	var buf bytes.Buffer
	b, _ := NewBudget(time.Minute, time.Second, &buf)

	exceeded := b.Record(100 * time.Millisecond)
	if exceeded {
		t.Fatal("expected budget not exceeded")
	}
	if buf.Len() > 0 {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}

func TestRecord_ExceedsBudget(t *testing.T) {
	var buf bytes.Buffer
	b, _ := NewBudget(time.Minute, 500*time.Millisecond, &buf)

	b.Record(300 * time.Millisecond)
	exceeded := b.Record(300 * time.Millisecond)
	if !exceeded {
		t.Fatal("expected budget exceeded")
	}
	if !strings.Contains(buf.String(), "exceeds threshold") {
		t.Fatalf("expected warning in output, got: %s", buf.String())
	}
}

func TestTotal_ReflectsCumulativeDuration(t *testing.T) {
	b, _ := NewBudget(time.Minute, 10*time.Second, nil)
	b.Record(200 * time.Millisecond)
	b.Record(300 * time.Millisecond)

	total := b.Total()
	if total < 500*time.Millisecond {
		t.Fatalf("expected total >= 500ms, got %s", total)
	}
}

func TestReset_ClearsDurations(t *testing.T) {
	b, _ := NewBudget(time.Minute, 10*time.Second, nil)
	b.Record(500 * time.Millisecond)
	b.Reset()

	if total := b.Total(); total != 0 {
		t.Fatalf("expected zero total after reset, got %s", total)
	}
}

func TestRecord_PrunesOldEntries(t *testing.T) {
	b, _ := NewBudget(50*time.Millisecond, 10*time.Second, nil)
	b.Record(200 * time.Millisecond)

	time.Sleep(60 * time.Millisecond)

	// old entry should be pruned; total should only include new record
	b.Record(100 * time.Millisecond)
	total := b.Total()
	if total > 150*time.Millisecond {
		t.Fatalf("expected old entry pruned, got total %s", total)
	}
}
