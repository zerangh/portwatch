package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllow_IndependentKeys(t *testing.T) {
	l := ratelimit.New(10 * time.Second)
	if !l.Allow("host-a") {
		t.Fatal("expected host-a first call to be allowed")
	}
	if !l.Allow("host-b") {
		t.Fatal("expected host-b first call to be allowed")
	}
}

func TestAllow_SameKeyThrottled(t *testing.T) {
	l := ratelimit.New(10 * time.Second)
	l.Allow("host-a")
	if l.Allow("host-a") {
		t.Fatal("expected second call for host-a to be throttled")
	}
}

func TestReset_UnblocksKey(t *testing.T) {
	l := ratelimit.New(10 * time.Second)
	l.Allow("host-a")
	l.Reset("host-a")
	if !l.Allow("host-a") {
		t.Fatal("expected call after reset to be allowed")
	}
}

func TestWait_ContextCancelled(t *testing.T) {
	l := ratelimit.New(10 * time.Second)
	l.Allow("host-a")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := l.Wait(ctx, "host-a"); err == nil {
		t.Fatal("expected error on cancelled context")
	}
}

func TestWait_AllowsAfterInterval(t *testing.T) {
	l := ratelimit.New(20 * time.Millisecond)
	l.Allow("host-x")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := l.Wait(ctx, "host-x"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
