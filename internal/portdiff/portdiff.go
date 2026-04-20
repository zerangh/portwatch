// Package portdiff provides utilities for computing and formatting
// the difference between two sets of open ports.
package portdiff

import (
	"fmt"
	"sort"
	"strings"
)

// Result holds the computed difference between two port sets.
type Result struct {
	Added   []int
	Removed []int
}

// IsEmpty returns true when there are no changes.
func (r Result) IsEmpty() bool {
	return len(r.Added) == 0 && len(r.Removed) == 0
}

// Summary returns a human-readable one-line description of the diff.
func (r Result) Summary() string {
	if r.IsEmpty() {
		return "no changes"
	}
	parts := make([]string, 0, 2)
	if len(r.Added) > 0 {
		parts = append(parts, fmt.Sprintf("+%d added", len(r.Added)))
	}
	if len(r.Removed) > 0 {
		parts = append(parts, fmt.Sprintf("-%d removed", len(r.Removed)))
	}
	return strings.Join(parts, ", ")
}

// Compute returns the diff between prev and next port slices.
// Both slices may be unsorted; the result slices are sorted ascending.
func Compute(prev, next []int) Result {
	prevSet := toSet(prev)
	nextSet := toSet(next)

	var added, removed []int
	for p := range nextSet {
		if !prevSet[p] {
			added = append(added, p)
		}
	}
	for p := range prevSet {
		if !nextSet[p] {
			removed = append(removed, p)
		}
	}
	sort.Ints(added)
	sort.Ints(removed)
	return Result{Added: added, Removed: removed}
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
