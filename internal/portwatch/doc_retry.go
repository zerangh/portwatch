// Package portwatch provides the core watch pipeline, runner, and supporting
// utilities for portwatch.
//
// # Retry
//
// RetryPolicy controls how transient scan failures are handled within a single
// scheduled run.  When a scan function returns an error the policy will sleep
// for Delay and try again up to MaxAttempts times before propagating the
// failure.  Context cancellation is respected between attempts so the process
// can be shut down cleanly without waiting for the full delay.
//
// Usage:
//
//	p := portwatch.DefaultRetryPolicy()
//	err := p.Retry(ctx, func() error {
//		return scanner.Scan()
//	})
//
package portwatch
