// Package history provides utilities for recording, querying, and managing
// port scan history over time.
//
// # Pruning
//
// The Prune function applies a RetentionPolicy to an existing history file,
// removing entries that exceed the configured maximum age or entry count.
//
// Example:
//
//	policy := history.DefaultRetentionPolicy()
//	policy.MaxAge = 7 * 24 * time.Hour
//	removed, err := history.Prune("/var/lib/portwatch/history.json", policy)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("pruned %d old entries\n", removed)
package history
