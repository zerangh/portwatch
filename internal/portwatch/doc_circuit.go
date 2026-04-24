// Package portwatch provides the core scan pipeline, runner, and supporting
// coordination primitives for portwatch.
//
// # Circuit Breaker
//
// Circuit implements a three-state circuit breaker (closed → open → half-open)
// that wraps the scan pipeline to prevent repeated failing scans from flooding
// logs or exhausting system resources.
//
// Usage:
//
//	cb, err := portwatch.NewCircuit(5, 30*time.Second, os.Stderr)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if cb.Allow() {
//		err := pipeline.Run(ctx)
//		if err != nil {
//			cb.RecordFailure()
//		} else {
//			cb.RecordSuccess()
//		}
//	}
//
// After maxFailures consecutive failures the circuit opens and blocks further
// attempts. Once resetAfter has elapsed a single probe is permitted; on
// success the circuit closes again.
package portwatch
