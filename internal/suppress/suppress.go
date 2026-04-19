// Package suppress provides a suppression list for ports that should never
// trigger alerts, regardless of their open/closed state.
package suppress

import (
	"encoding/json"
	"os"
	"sync"
)

// List holds a set of suppressed ports.
type List struct {
	mu    sync.RWMutex
	ports map[int]struct{}
	path  string
}

// New returns an empty suppression list backed by path.
func New(path string) (*List, error) {
	if path == "" {
		return nil, fmt.Errorf("suppress: path must not be empty")
	}
	l := &List{path: path, ports: make(map[int]struct{})}
	if err := l.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return l, nil
}

func (l *List) load() error {
	data, err := os.ReadFile(l.path)
	if err != nil {
		return err
	}
	var ports []int
	if err := json.Unmarshal(data, &ports); err != nil {
		return fmt.Errorf("suppress: parse %s: %w", l.path, err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ports = make(map[int]struct{}, len(ports))
	for _, p := range ports {
		l.ports[p] = struct{}{}
	}
	return nil
}

// Add adds a port to the suppression list and persists it.
func (l *List) Add(port int) error {
	l.mu.Lock()
	l.ports[port] = struct{}{}
	l.mu.Unlock()
	return l.save()
}

// Remove removes a port from the suppression list and persists it.
func (l *List) Remove(port int) error {
	l.mu.Lock()
	delete(l.ports, port)
	l.mu.Unlock()
	return l.save()
}

// Contains reports whether port is suppressed.
func (l *List) Contains(port int) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, ok := l.ports[port]
	return ok
}

// Filter returns only those ports not present in the suppression list.
func (l *List) Filter(ports []int) []int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := ports[:0:0]
	for _, p := range ports {
		if _, suppressed := l.ports[p]; !suppressed {
			out = append(out, p)
		}
	}
	return out
}

func (l *List) save() error {
	l.mu.RLock()
	ports := make([]int, 0, len(l.ports))
	for p := range l.ports {
		ports = append(ports, p)
	}
	l.mu.RUnlock()
	data, err := json.MarshalIndent(ports, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(l.path, data, 0o644)
}
