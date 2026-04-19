// Package watchdog monitors scan cycles and raises alerts when scans
// have not completed within an expected interval.
package watchdog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Watchdog tracks the last successful scan time and reports staleness.
type Watchdog struct {
	mu       sync.Mutex
	last     time.Time
	maxAge   time.Duration
	writer   io.Writer
}

// New creates a Watchdog that considers scans stale after maxAge.
func New(maxAge time.Duration, w io.Writer) (*Watchdog, error) {
	if maxAge <= 0 {
		return nil, fmt.Errorf("watchdog: maxAge must be positive")
	}
	if w == nil {
		w = os.Stderr
	}
	return &Watchdog{maxAge: maxAge, writer: w}, nil
}

// Ping records a successful scan at the current time.
func (wd *Watchdog) Ping() {
	wd.mu.Lock()
	defer wd.mu.Unlock()
	wd.last = time.Now()
}

// IsStale returns true when no ping has been received within maxAge.
func (wd *Watchdog) IsStale() bool {
	wd.mu.Lock()
	defer wd.mu.Unlock()
	if wd.last.IsZero() {
		return true
	}
	return time.Since(wd.last) > wd.maxAge
}

// Age returns the duration since the last ping, or -1 if never pinged.
func (wd *Watchdog) Age() time.Duration {
	wd.mu.Lock()
	defer wd.mu.Unlock()
	if wd.last.IsZero() {
		return -1
	}
	return time.Since(wd.last)
}

// Check writes a warning to the writer if the watchdog is stale.
func (wd *Watchdog) Check() {
	if wd.IsStale() {
		fmt.Fprintf(wd.writer, "[watchdog] WARNING: no scan completed in the last %s\n", wd.maxAge)
	}
}
