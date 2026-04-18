package history

import "time"

// Query holds filter parameters for searching history entries.
type Query struct {
	Since  time.Time
	Until  time.Time
	Limit  int
	HasDiff bool // only entries with changes
}

// Filter returns entries from h that match q.
func (h *History) Filter(q Query) []Entry {
	var out []Entry
	for _, e := range h.entries {
		if !q.Since.IsZero() && e.Timestamp.Before(q.Since) {
			continue
		}
		if !q.Until.IsZero() && e.Timestamp.After(q.Until) {
			continue
		}
		if q.HasDiff && len(e.Added) == 0 && len(e.Removed) == 0 {
			continue
		}
		out = append(out, e)
		if q.Limit > 0 && len(out) >= q.Limit {
			break
		}
	}
	return out
}

// Since returns all entries recorded after t.
func (h *History) Since(t time.Time) []Entry {
	return h.Filter(Query{Since: t})
}

// WithChanges returns only entries that recorded port changes.
func (h *History) WithChanges() []Entry {
	return h.Filter(Query{HasDiff: true})
}
