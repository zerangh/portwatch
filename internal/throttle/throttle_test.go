package throttle_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestAllow_FirstCallPermitted(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	th.Allow()
	if th.Allow() {
		t.Fatal("expected second immediate call to be blocked")
	}
}

func TestAllow_AfterInterval(t *testing.T) {
	th := throttle.New(20 * time.Millisecond)
	th.Allow()
	time.Sleep(30 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	th := throttle.New(10 * time.Second)
	th.Allow()
	th.Reset()
	if !th.Allow() {
		t.Fatal("expected call after reset to be allowed")
	}
}

func TestWait_CompletesWithinInterval(t *testing.T) {
	th := throttle.New(20 * time.Millisecond)
	th.Allow()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	start := time.Now()
	if err := th.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if time.Since(start) < 15*time.Millisecond {
		t.Fatal("wait returned too quickly")
	}
}

func TestWait_CancelledContext(t *testing.T) {
	th := throttle.New(10 * time.Second)
	th.Allow()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := th.Wait(ctx); err == nil {
		t.Fatal("expected error on cancelled context")
	}
}
