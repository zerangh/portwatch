// Package history provides scan history recording, querying, and analysis.
//
// # Summary
//
// The Summarize function aggregates a slice of history entries into a Summary
// struct, providing at-a-glance statistics:
//
//	- Total: number of recorded scans
//	- WithChanges: scans that detected at least one port change
//	- FirstSeen / LastSeen: time range covered by the entries
//	- MostAdded / MostRemoved: peak port-change counts in a single scan
//
// Use ChangeRate() on a Summary to get the fraction of scans that had changes.
//
//	entries, _ := h.Load()
//	s := history.Summarize(entries)
//	fmt.Printf("Change rate: %.1f%%\n", s.ChangeRate()*100)
package history
