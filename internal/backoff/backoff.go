// Package backoff provides exponential backoff retry logic for portwatch operations.
package backoff

import (
	"context"
	"math"
	"time"
)

// Policy defines the backoff configuration.
type Policy struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	MaxAttempts     int
}

// DefaultPolicy returns a sensible default backoff policy.
func DefaultPolicy() Policy {
	return Policy{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		MaxAttempts:     5,
	}
}

// Backoff holds state for a single retry sequence.
type Backoff struct {
	policy  Policy
	attempt int
}

// New creates a new Backoff from the given policy.
func New(p Policy) *Backoff {
	return &Backoff{policy: p}
}

// Next returns the duration to wait before the next attempt and whether retrying should continue.
func (b *Backoff) Next() (time.Duration, bool) {
	if b.policy.MaxAttempts > 0 && b.attempt >= b.policy.MaxAttempts {
		return 0, false
	}
	interval := float64(b.policy.InitialInterval) * math.Pow(b.policy.Multiplier, float64(b.attempt))
	if interval > float64(b.policy.MaxInterval) {
		interval = float64(b.policy.MaxInterval)
	}
	b.attempt++
	return time.Duration(interval), true
}

// Reset restarts the backoff sequence.
func (b *Backoff) Reset() {
	b.attempt = 0
}

// Attempt returns the current attempt count.
func (b *Backoff) Attempt() int {
	return b.attempt
}

// Wait sleeps for the next backoff interval, respecting context cancellation.
// Returns false if retries are exhausted or context is done.
func (b *Backoff) Wait(ctx context.Context) bool {
	d, ok := b.Next()
	if !ok {
		return false
	}
	select {
	case <-time.After(d):
		return true
	case <-ctx.Done():
		return false
	}
}
