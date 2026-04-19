package audit

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// ScanEvent holds metadata about a single scan cycle.
type ScanEvent struct {
	Host     string
	Duration time.Duration
	Ports    int
	Err      error
}

// RecordScan logs a completed scan event, including any error.
func (l *Logger) RecordScan(e ScanEvent) {
	if e.Err != nil {
		_ = l.Error("scan.error", fmt.Sprintf("host=%s err=%v", e.Host, e.Err))
		return
	}
	_ = l.Info("scan.complete", fmt.Sprintf("host=%s ports=%d duration=%s",
		e.Host, e.Ports, e.Duration.Round(time.Millisecond)))
}

// RecordDiff logs added and removed ports derived from a snapshot comparison.
func (l *Logger) RecordDiff(host string, diff snapshot.Comparison) {
	if len(diff.Added) == 0 && len(diff.Removed) == 0 {
		return
	}
	if len(diff.Added) > 0 {
		_ = l.Warn("ports.added", fmt.Sprintf("host=%s ports=%v", host, diff.Added))
	}
	if len(diff.Removed) > 0 {
		_ = l.Warn("ports.removed", fmt.Sprintf("host=%s ports=%v", host, diff.Removed))
	}
}
