// Package history tracks and persists the results of port scans over time.
//
// Each scan produces an Entry containing the timestamp, the full list of
// open ports, and the ports that were added or removed relative to the
// previous scan. Entries are appended to a JSON file on disk so that
// operators can review trends and audit changes across restarts.
//
// Usage:
//
//	h := history.New("/var/lib/portwatch/history.json")
//	if err := h.Load(); err != nil { ... }
//	h.Append(history.Entry{...})
package history
