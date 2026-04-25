package portwatch

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Quota enforces a maximum number of scan cycles within a rolling time window.
// Once the quota is exhausted the caller must wait until the window resets.
type Quota struct {
	mu      sync.Mutex
	max     int
	window  time.Duration
	count   int
	windowStart time.Time
	w       io.Writer
}

// NewQuota creates a Quota that allows at most max scans per window duration.
// A nil writer falls back to os.Stderr.
func NewQuota(max int, window time.Duration, w io.Writer) (*Quota, error) {
	if max <= 0 {
		return nil, fmt.Errorf("quota: max must be > 0, got %d", max)
	}
	if window <= 0 {
		return nil, fmt.Errorf("quota: window must be > 0, got %s", window)
	}
	if w == nil {
		w = os.Stderr
	}
	return &Quota{
		max:         max,
		window:      window,
		windowStart: time.Now(),
		w:           w,
	}, nil
}

// Allow reports whether a scan is permitted under the current quota.
// It advances the internal counter and resets the window when expired.
func (q *Quota) Allow() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	if now.Sub(q.windowStart) >= q.window {
		q.count = 0
		q.windowStart = now
	}

	if q.count >= q.max {
		fmt.Fprintf(q.w, "quota: limit of %d scans per %s reached, skipping\n", q.max, q.window)
		return false
	}

	q.count++
	return true
}

// Remaining returns the number of scans still available in the current window.
func (q *Quota) Remaining() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	if time.Since(q.windowStart) >= q.window {
		return q.max
	}
	remaining := q.max - q.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset clears the counter and starts a fresh window immediately.
func (q *Quota) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.count = 0
	q.windowStart = time.Now()
}
