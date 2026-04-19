package debounce_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
)

func TestCall_FiresAfterDelay(t *testing.T) {
	d := debounce.New(50 * time.Millisecond)
	var count int32
	d.Call(func() { atomic.AddInt32(&count, 1) })
	time.Sleep(100 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Fatalf("expected 1 call, got %d", got)
	}
}

func TestCall_CoalescesRapidCalls(t *testing.T) {
	d := debounce.New(60 * time.Millisecond)
	var count int32
	for i := 0; i < 5; i++ {
		d.Call(func() { atomic.AddInt32(&count, 1) })
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(120 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Fatalf("expected 1 coalesced call, got %d", got)
	}
}

func TestPending_TrueBeforeFire(t *testing.T) {
	d := debounce.New(100 * time.Millisecond)
	d.Call(func() {})
	if !d.Pending() {
		t.Fatal("expected pending to be true before delay")
	}
	time.Sleep(150 * time.Millisecond)
	if d.Pending() {
		t.Fatal("expected pending to be false after delay")
	}
}

func TestFlush_CancelsPending(t *testing.T) {
	d := debounce.New(200 * time.Millisecond)
	var fired int32
	d.Call(func() { atomic.AddInt32(&fired, 1) })
	if !d.Flush() {
		t.Fatal("expected Flush to return true")
	}
	time.Sleep(250 * time.Millisecond)
	if got := atomic.LoadInt32(&fired); got != 0 {
		t.Fatalf("expected 0 calls after flush, got %d", got)
	}
}

func TestFlush_NothingPending(t *testing.T) {
	d := debounce.New(50 * time.Millisecond)
	if d.Flush() {
		t.Fatal("expected Flush to return false when nothing pending")
	}
}
