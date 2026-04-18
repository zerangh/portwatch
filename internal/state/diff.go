package state

import "sort"

// Diff holds the result of comparing two port states.
type Diff struct {
	Added   []int
	Removed []int
}

// IsEmpty returns true when there are no changes.
func (d Diff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0
}

// Summary returns a human-readable one-line summary of the diff.
func (d Diff) Summary() string {
	if d.IsEmpty() {
		return "no changes"
	}
	return formatSummary(d.Added, d.Removed)
}

func formatSummary(added, removed []int) string {
	parts := ""
	if len(added) > 0 {
		sort.Ints(added)
		parts += formatPart("opened", added)
	}
	if len(removed) > 0 {
		sort.Ints(removed)
		if parts != "" {
			parts += "; "
		}
		parts += formatPart("closed", removed)
	}
	return parts
}

func formatPart(label string, ports []int) string {
	s := label + ": "
	for i, p := range ports {
		if i > 0 {
			s += ","
		}
		s += itoa(p)
	}
	return s
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}
