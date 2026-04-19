// Package throttle provides rate-limiting utilities for port scan cycles.
package throttle

import (
	"context"
	"sync"
	"time"
)

// Throttle limits how frequently an action can be triggered.
type Throttle struct {
	mu       sync.Mutex
	interval time.Duration
	last     time.Time
}

// New creates a Throttle that allows at most one action per interval.
func New(interval time.Duration) *Throttle {
	return &Throttle{interval: interval}
}

// Allow returns true if enough time has elapsed since the last allowed call.
func (t *Throttle) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	if now.Sub(t.last) >= t.interval {
		t.last = now
		return true
	}
	return false
}

// Wait blocks until the throttle allows the next action or ctx is cancelled.
func (t *Throttle) Wait(ctx context.Context) error {
	for {
		if t.Allow() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// Reset clears the last-allowed timestamp, immediately permitting the next call.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = time.Time{}
}
