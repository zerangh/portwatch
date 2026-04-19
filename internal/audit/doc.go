// Package audit provides structured audit logging for portwatch events.
//
// An audit Logger records timestamped entries for significant events such as
// scan starts, completions, port changes, and errors. Output can be formatted
// as human-readable text or newline-delimited JSON for ingestion by log
// aggregators.
//
// Usage:
//
//	l := audit.New(os.Stderr, false)
//	l.Info("scan.start", "range 1-1024")
//	l.Warn("port.changed", "added: 8080, 9090")
package audit
