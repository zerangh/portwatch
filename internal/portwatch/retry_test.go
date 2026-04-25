package portwatch

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

var errTest = errors.New("test error")

func TestDefaultRetryPolicy_Values(t *testing.T) {
	p := DefaultRetryPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
	if p.Delay != 2*time.Second {
		t.Errorf("expected Delay=2s, got %v", p.Delay)
	}
	if p.Writer == nil {
		t.Error("expected non-nil Writer")
	}
}

func TestRetry_SucceedsFirstAttempt(t *testing.T) {
	p := RetryPolicy{MaxAttempts: 3, Delay: 0, Writer: &bytes.Buffer{}}
	calls := 0
	err := p.Retry(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetry_RetriesOnFailure(t *testing.T) {
	var buf bytes.Buffer
	p := RetryPolicy{MaxAttempts: 3, Delay: 0, Writer: &buf}
	calls := 0
	err := p.Retry(context.Background(), func() error {
		calls++
		if calls < 3 {
			return errTest
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetry_ExhaustsAttempts(t *testing.T) {
	var buf bytes.Buffer
	p := RetryPolicy{MaxAttempts: 2, Delay: 0, Writer: &buf}
	err := p.Retry(context.Background(), func() error { return errTest })
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "2 attempts") {
		t.Errorf("error should mention attempt count: %v", err)
	}
	if !errors.Is(err, errTest) {
		t.Errorf("expected wrapped errTest, got %v", err)
	}
}

func TestRetry_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := RetryPolicy{MaxAttempts: 5, Delay: time.Second, Writer: &bytes.Buffer{}}
	err := p.Retry(ctx, func() error { return errTest })
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestRetry_ZeroMaxAttemptsBecomesOne(t *testing.T) {
	p := RetryPolicy{MaxAttempts: 0, Delay: 0, Writer: &bytes.Buffer{}}
	calls := 0
	_ = p.Retry(context.Background(), func() error {
		calls++
		return errTest
	})
	if calls != 1 {
		t.Errorf("expected 1 call for zero MaxAttempts, got %d", calls)
	}
}

func TestRetry_NilWriterUsesStderr(t *testing.T) {
	p := RetryPolicy{MaxAttempts: 1, Delay: 0, Writer: nil}
	err := p.Retry(context.Background(), func() error { return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
