// Package throttle provides a simple rate-limiter used to prevent portwatch
// from triggering scan cycles more frequently than a configured minimum
// interval.
//
// Usage:
//
//	th := throttle.New(30 * time.Second)
//
//	// Non-blocking check:
//	if th.Allow() {
//		// run scan
//	}
//
//	// Blocking wait:
//	if err := th.Wait(ctx); err != nil {
//		// context cancelled
//	}
package throttle
