// Package snapshot provides point-in-time capture of open ports
// and comparison utilities for detecting changes between snapshots.
package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// Snapshot represents a captured set of open ports at a point in time.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
}

// New creates a new Snapshot with the current timestamp.
func New(host string, ports []int) *Snapshot {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Host:      host,
		Ports:     sorted,
	}
}

// Contains reports whether the snapshot includes the given port.
func (s *Snapshot) Contains(port int) bool {
	for _, p := range s.Ports {
		if p == port {
			return true
		}
	}
	return false
}

// Compare returns ports added and removed relative to a previous snapshot.
func (s *Snapshot) Compare(prev *Snapshot) (added, removed []int) {
	if prev == nil {
		return s.Ports, nil
	}
	prevSet := toSet(prev.Ports)
	currSet := toSet(s.Ports)

	for p := range currSet {
		if !prevSet[p] {
			added = append(added, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			removed = append(removed, p)
		}
	}
	sort.Ints(added)
	sort.Ints(removed)
	return
}

// String returns a human-readable summary of the snapshot.
func (s *Snapshot) String() string {
	return fmt.Sprintf("Snapshot{host=%s, ports=%d, at=%s}",
		s.Host, len(s.Ports), s.Timestamp.Format(time.RFC3339))
}

func toSet(ports []int) map[int]bool {
	m := make(map[int]bool, len(ports))
	for _, p := range ports {
		m[p] = true
	}
	return m
}
