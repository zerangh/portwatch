package portwatch

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Budget tracks cumulative scan durations within a rolling time window and
// warns when the total scanning time exceeds a configured threshold.
type Budget struct {
	mu        sync.Mutex
	window    time.Duration
	threshold time.Duration
	entries   []budgetEntry
	w         io.Writer
}

type budgetEntry struct {
	at       time.Time
	duration time.Duration
}

// NewBudget returns a Budget that warns when cumulative scan time within
// window exceeds threshold. A nil writer falls back to os.Stderr.
func NewBudget(window, threshold time.Duration, w io.Writer) (*Budget, error) {
	if window <= 0 {
		return nil, fmt.Errorf("budget: window must be positive")
	}
	if threshold <= 0 {
		return nil, fmt.Errorf("budget: threshold must be positive")
	}
	if w == nil {
		w = os.Stderr
	}
	return &Budget{window: window, threshold: threshold, w: w}, nil
}

// Record adds a scan duration and returns true if the cumulative budget
// within the window has been exceeded.
func (b *Budget) Record(d time.Duration) bool {
	now := time.Now()
	b.mu.Lock()
	defer b.mu.Unlock()

	b.prune(now)
	b.entries = append(b.entries, budgetEntry{at: now, duration: d})

	total := b.total()
	if total > b.threshold {
		fmt.Fprintf(b.w, "[budget] cumulative scan time %s exceeds threshold %s in window %s\n",
			total.Round(time.Millisecond), b.threshold, b.window)
		return true
	}
	return false
}

// Total returns the cumulative scan duration within the current window.
func (b *Budget) Total() time.Duration {
	now := time.Now()
	b.mu.Lock()
	defer b.mu.Unlock()
	b.prune(now)
	return b.total()
}

// Reset clears all recorded entries.
func (b *Budget) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = b.entries[:0]
}

func (b *Budget) prune(now time.Time) {
	cutoff := now.Add(-b.window)
	i := 0
	for i < len(b.entries) && b.entries[i].at.Before(cutoff) {
		i++
	}
	b.entries = b.entries[i:]
}

func (b *Budget) total() time.Duration {
	var sum time.Duration
	for _, e := range b.entries {
		sum += e.duration
	}
	return sum
}
