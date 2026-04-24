package portwatch

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Circuit is a simple circuit breaker for scan pipelines.
type Circuit struct {
	mu           sync.Mutex
	state        State
	failures      int
	maxFailures   int
	resetAfter    time.Duration
	openedAt      time.Time
	w             io.Writer
}

// NewCircuit creates a Circuit breaker that opens after maxFailures consecutive
// errors and attempts recovery after resetAfter.
func NewCircuit(maxFailures int, resetAfter time.Duration, w io.Writer) (*Circuit, error) {
	if maxFailures <= 0 {
		return nil, fmt.Errorf("circuit: maxFailures must be > 0")
	}
	if resetAfter <= 0 {
		return nil, fmt.Errorf("circuit: resetAfter must be > 0")
	}
	if w == nil {
		w = os.Stderr
	}
	return &Circuit{
		state:       StateClosed,
		maxFailures: maxFailures,
		resetAfter:  resetAfter,
		w:           w,
	}, nil
}

// Allow returns true if the circuit permits a scan attempt.
func (c *Circuit) Allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch c.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(c.openedAt) >= c.resetAfter {
			c.state = StateHalfOpen
			fmt.Fprintf(c.w, "circuit: half-open, probing\n")
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

// RecordSuccess resets the circuit to closed.
func (c *Circuit) RecordSuccess() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.state != StateClosed {
		fmt.Fprintf(c.w, "circuit: closed after recovery\n")
	}
	c.failures = 0
	c.state = StateClosed
}

// RecordFailure increments the failure count and may open the circuit.
func (c *Circuit) RecordFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures++
	if c.failures >= c.maxFailures && c.state != StateOpen {
		c.state = StateOpen
		c.openedAt = time.Now()
		fmt.Fprintf(c.w, "circuit: opened after %d failures\n", c.failures)
	}
}

// State returns the current circuit state.
func (c *Circuit) State() State {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}
