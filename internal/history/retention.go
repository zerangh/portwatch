package history

import (
	"sort"
	"time"
)

// RetentionPolicy defines how long history entries are kept.
type RetentionPolicy struct {
	// MaxAge is the maximum age of an entry before it is pruned.
	// Zero means no age-based pruning.
	MaxAge time.Duration

	// MaxEntries is the maximum number of entries to retain.
	// Zero means no count-based pruning.
	MaxEntries int
}

// DefaultRetentionPolicy returns a sensible default: keep last 100 entries
// and discard anything older than 30 days.
func DefaultRetentionPolicy() RetentionPolicy {
	return RetentionPolicy{
		MaxAge:     30 * 24 * time.Hour,
		MaxEntries: 100,
	}
}

// Apply prunes entries according to the policy and returns the surviving
// entries sorted newest-first.
func (p RetentionPolicy) Apply(entries []Entry) []Entry {
	if len(entries) == 0 {
		return entries
	}

	// Sort newest first.
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.After(sorted[j].Timestamp)
	})

	var kept []Entry
	cutoff := time.Now().Add(-p.MaxAge)

	for _, e := range sorted {
		if p.MaxAge > 0 && e.Timestamp.Before(cutoff) {
			continue
		}
		kept = append(kept, e)
		if p.MaxEntries > 0 && len(kept) >= p.MaxEntries {
			break
		}
	}

	if kept == nil {
		return []Entry{}
	}
	return kept
}
