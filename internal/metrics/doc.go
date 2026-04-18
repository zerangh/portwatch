// Package metrics provides a lightweight, thread-safe counter for tracking
// portwatch scan statistics at runtime.
//
// Usage:
//
//	ctr := metrics.New()
//
//	// after each scan:
//	ctr.Record(changed, err)
//
//	// inspect current totals:
//	snap := ctr.Snapshot()
//	fmt.Printf("scans=%d changes=%d errors=%d\n",
//		snap.Scans, snap.Changes, snap.Errors)
//
// Snapshot returns a value copy so callers can inspect counters without
// holding the internal lock.
package metrics
