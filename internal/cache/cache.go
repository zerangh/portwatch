// Package cache provides a simple in-memory cache for port scan results
// to avoid redundant scans within a configurable TTL window.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached scan result.
type Entry struct {
	Ports     []int
	CachedAt  time.Time
	ExpiresAt time.Time
}

// Cache stores scan results keyed by host.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	ttl     time.Duration
}

// New creates a new Cache with the given TTL.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]*Entry),
		ttl:     ttl,
	}
}

// Set stores ports for the given host key.
func (c *Cache) Set(key string, ports []int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.entries[key] = &Entry{
		Ports:     ports,
		CachedAt:  now,
		ExpiresAt: now.Add(c.ttl),
	}
}

// Get returns the cached ports and true if a valid (non-expired) entry exists.
func (c *Cache) Get(key string) ([]int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok || time.Now().After(e.ExpiresAt) {
		return nil, false
	}
	return e.Ports, true
}

// Invalidate removes the entry for the given key.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*Entry)
}
