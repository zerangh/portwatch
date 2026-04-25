package portwatch

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// RetryPolicy controls how scan failures are retried within a single run.
type RetryPolicy struct {
	MaxAttempts int
	Delay       time.Duration
	Writer      io.Writer
}

// DefaultRetryPolicy returns a RetryPolicy with sensible defaults.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts: 3,
		Delay:       2 * time.Second,
		Writer:      os.Stderr,
	}
}

// Retry executes fn up to MaxAttempts times, returning the first nil error.
// Each failed attempt (except the last) sleeps for Delay before retrying.
// If ctx is cancelled the function returns ctx.Err() immediately.
func (p RetryPolicy) Retry(ctx context.Context, fn func() error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	w := p.Writer
	if w == nil {
		w = os.Stderr
	}

	var last error
	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		last = fn()
		if last == nil {
			return nil
		}
		fmt.Fprintf(w, "portwatch: scan attempt %d/%d failed: %v\n", attempt, p.MaxAttempts, last)
		if attempt < p.MaxAttempts {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(p.Delay):
			}
		}
	}
	return fmt.Errorf("scan failed after %d attempts: %w", p.MaxAttempts, last)
}
