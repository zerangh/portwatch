// Package metrics tracks runtime statistics for portwatch scans.
package metrics

import (
	"sync"
	"time"
)

// Counter holds cumulative scan statistics.
type Counter struct {
	mu          sync.Mutex
	Scans       int
	Changes     int
	Errors      int
	LastScan    time.Time
	LastChange  time.Time
}

// Record registers the result of a completed scan.
func (c *Counter) Record(changed bool, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Scans++
	c.LastScan = time.Now()
	if err != nil {
		c.Errors++
		return
	}
	if changed {
		c.Changes++
		c.LastChange = time.Now()
	}
}

// Snapshot returns a copy of the current counters.
func (c *Counter) Snapshot() Counter {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Counter{
		Scans:      c.Scans,
		Changes:    c.Changes,
		Errors:     c.Errors,
		LastScan:   c.LastScan,
		LastChange: c.LastChange,
	}
}

// New returns an initialised Counter.
func New() *Counter {
	return &Counter{}
}
