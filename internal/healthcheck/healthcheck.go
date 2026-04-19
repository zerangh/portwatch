// Package healthcheck provides a simple health status reporter for portwatch.
package healthcheck

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Status represents the health state of the watcher.
type Status struct {
	Healthy     bool      `json:"healthy"`
	LastScan    time.Time `json:"last_scan"`
	LastError   string    `json:"last_error,omitempty"`
	ScanCount   int       `json:"scan_count"`
	ErrorCount  int       `json:"error_count"`
}

// Checker tracks scan outcomes and reports health.
type Checker struct {
	status Status
	w      io.Writer
}

// New returns a new Checker writing to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Checker {
	if w == nil {
		w = os.Stdout
	}
	return &Checker{w: w}
}

// RecordScan records a successful scan at the given time.
func (c *Checker) RecordScan(t time.Time) {
	c.status.LastScan = t
	c.status.ScanCount++
	c.status.Healthy = true
}

// RecordError records a failed scan.
func (c *Checker) RecordError(err error) {
	c.status.ErrorCount++
	c.status.Healthy = false
	if err != nil {
		c.status.LastError = err.Error()
	}
}

// Status returns a copy of the current health status.
func (c *Checker) Status() Status {
	return c.status
}

// Print writes a human-readable health summary to the writer.
func (c *Checker) Print() {
	s := c.status
	healthStr := "OK"
	if !s.Healthy {
		healthStr = "DEGRADED"
	}
	fmt.Fprintf(c.w, "health: %s | scans: %d | errors: %d | last_scan: %s\n",
		healthStr, s.ScanCount, s.ErrorCount, formatTime(s.LastScan))
	if s.LastError != "" {
		fmt.Fprintf(c.w, "last_error: %s\n", s.LastError)
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	return t.Format(time.RFC3339)
}
