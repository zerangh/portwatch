// Package alert handles user-facing notifications when port state changes
// are detected between two consecutive scans.
//
// Usage:
//
//	n := alert.New(os.Stdout)
//	if n.Notify(diff) {
//		// changes were printed
//	}
//
// The Notifier is intentionally simple — it writes plain-text output
// suitable for terminals and log aggregators. Future backends (e.g.
// webhook, email) can be added by extending this package.
package alert
