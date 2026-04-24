package portwatch

import (
	"sync"
	"time"
)

// Cooldown prevents repeated alert firing for the same port change within a
// configurable window. Once a change is recorded, subsequent identical changes
// are suppressed until the cooldown period expires.
type Cooldown struct {
	mu       sync.Mutex
	window   time.Duration
	recorded map[string]time.Time
	now      func() time.Time
}

// NewCooldown returns a Cooldown with the given suppression window.
// A zero or negative window disables suppression (all events pass through).
func NewCooldown(window time.Duration) *Cooldown {
	return &Cooldown{
		window:   window,
		recorded: make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow reports whether the event identified by key should be allowed through.
// If the key was seen within the cooldown window it is suppressed (false).
// Otherwise the key is recorded and true is returned.
func (c *Cooldown) Allow(key string) bool {
	if c.window <= 0 {
		return true
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	if last, ok := c.recorded[key]; ok && now.Sub(last) < c.window {
		return false
	}
	c.recorded[key] = now
	return true
}

// Reset clears the recorded timestamp for key, immediately allowing the next
// event with that key through regardless of the remaining cooldown window.
func (c *Cooldown) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.recorded, key)
}

// Flush removes all recorded entries whose cooldown window has expired,
// freeing memory for long-running processes.
func (c *Cooldown) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	for k, t := range c.recorded {
		if now.Sub(t) >= c.window {
			delete(c.recorded, k)
		}
	}
}

// Len returns the number of keys currently tracked (including those whose
// window may have expired but Flush has not yet been called).
func (c *Cooldown) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.recorded)
}
