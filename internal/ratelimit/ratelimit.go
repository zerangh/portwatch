// Package ratelimit wraps throttle to provide per-host scan rate limiting.
package ratelimit

import (
	"context"
	"sync"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

// Limiter manages per-key throttles.
type Limiter struct {
	mu       sync.Mutex
	interval time.Duration
	keys     map[string]*throttle.Throttle
}

// New creates a Limiter where each key is throttled to interval.
func New(interval time.Duration) *Limiter {
	return &Limiter{
		interval: interval,
		keys:     make(map[string]*throttle.Throttle),
	}
}

func (l *Limiter) get(key string) *throttle.Throttle {
	l.mu.Lock()
	defer l.mu.Unlock()
	if th, ok := l.keys[key]; ok {
		return th
	}
	th := throttle.New(l.interval)
	l.keys[key] = th
	return th
}

// Allow returns true if the key is not throttled.
func (l *Limiter) Allow(key string) bool {
	return l.get(key).Allow()
}

// Wait blocks until the key's throttle permits or ctx is cancelled.
func (l *Limiter) Wait(ctx context.Context, key string) error {
	return l.get(key).Wait(ctx)
}

// Reset clears the throttle state for a key.
func (l *Limiter) Reset(key string) {
	l.get(key).Reset()
}
