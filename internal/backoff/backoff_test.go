package backoff_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/backoff"
)

func TestDefaultPolicy(t *testing.T) {
	p := backoff.DefaultPolicy()
	if p.InitialInterval <= 0 {
		t.Fatal("expected positive initial interval")
	}
	if p.MaxInterval < p.InitialInterval {
		t.Fatal("max interval should be >= initial interval")
	}
	if p.Multiplier <= 1 {
		t.Fatal("multiplier should be > 1")
	}
}

func TestNext_Increases(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		MaxAttempts:     5,
	}
	b := backoff.New(p)
	prev, _ := b.Next()
	for i := 1; i < 4; i++ {
		d, ok := b.Next()
		if !ok {
			t.Fatalf("expected ok on attempt %d", i)
		}
		if d <= prev {
			t.Fatalf("expected interval to grow, got %v <= %v", d, prev)
		}
		prev = d
	}
}

func TestNext_RespectsMaxAttempts(t *testing.T) {
	p := backoff.Policy{InitialInterval: 10 * time.Millisecond, Multiplier: 2.0, MaxAttempts: 3}
	b := backoff.New(p)
	for i := 0; i < 3; i++ {
		_, ok := b.Next()
		if !ok {
			t.Fatalf("expected ok on attempt %d", i)
		}
	}
	_, ok := b.Next()
	if ok {
		t.Fatal("expected no more attempts after max")
	}
}

func TestNext_CapsAtMaxInterval(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 1 * time.Second,
		MaxInterval:     2 * time.Second,
		Multiplier:      10.0,
		MaxAttempts:     5,
	}
	b := backoff.New(p)
	for i := 0; i < 5; i++ {
		d, _ := b.Next()
		if d > 2*time.Second {
			t.Fatalf("interval %v exceeded max", d)
		}
	}
}

func TestReset_RestartsCounting(t *testing.T) {
	p := backoff.Policy{InitialInterval: 50 * time.Millisecond, Multiplier: 2.0, MaxAttempts: 2}
	b := backoff.New(p)
	b.Next()
	b.Next()
	b.Reset()
	if b.Attempt() != 0 {
		t.Fatalf("expected attempt 0 after reset, got %d", b.Attempt())
	}
	_, ok := b.Next()
	if !ok {
		t.Fatal("expected ok after reset")
	}
}

func TestWait_ContextCancelled(t *testing.T) {
	p := backoff.Policy{InitialInterval: 10 * time.Second, Multiplier: 1.0, MaxAttempts: 3}
	b := backoff.New(p)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if b.Wait(ctx) {
		t.Fatal("expected false when context already cancelled")
	}
}
