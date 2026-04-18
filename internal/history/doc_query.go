// Package history (query) provides filtering helpers for querying recorded
// scan history.
//
// Usage:
//
//	h, _ := history.New(path)
//
//	// All entries in the last 24 hours:
//	recent := h.Since(time.Now().Add(-24 * time.Hour))
//
//	// Only scans that detected a change:
//	changed := h.WithChanges()
//
//	// Flexible filtering:
//	results := h.Filter(history.Query{
//		Since:   time.Now().Add(-7 * 24 * time.Hour),
//		Limit:   10,
//		HasDiff: true,
//	})
package history
