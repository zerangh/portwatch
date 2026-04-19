// Package debounce provides a mechanism to delay and coalesce rapid
// successive calls into a single invocation after a quiet period.
package debounce

import (
	"sync"
	"time"
)

// Debouncer delays execution of a function until after a quiet period.
type Debouncer struct {
	mu       sync.Mutex
	delay    time.Duration
	timer    *time.Timer
	pending  bool
}

// New creates a new Debouncer with the given delay.
func New(delay time.Duration) *Debouncer {
	return &Debouncer{delay: delay}
}

// Call schedules fn to be called after the debounce delay.
// If Call is invoked again before the delay expires, the timer resets.
func (d *Debouncer) Call(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}

	d.pending = true
	d.timer = time.AfterFunc(d.delay, func() {
		d.mu.Lock()
		d.pending = false
		d.mu.Unlock()
		fn()
	})
}

// Pending reports whether a call is waiting to fire.
func (d *Debouncer) Pending() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.pending
}

// Flush cancels any pending timer and returns whether one was pending.
func (d *Debouncer) Flush() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil && d.pending {
		d.timer.Stop()
		d.pending = false
		return true
	}
	return false
}
