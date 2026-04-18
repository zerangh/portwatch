// Package state manages persistence of port scan snapshots and computes
// the difference between consecutive scans to detect port changes.
//
// Typical usage:
//
//	snap := state.Snapshot{
//		Timestamp: time.Now(),
//		Ports:     []int{22, 80, 443},
//	}
//	if err := state.Save("/var/lib/portwatch/state.json", snap); err != nil {
//		log.Fatal(err)
//	}
//
//	prev, err := state.Load("/var/lib/portwatch/state.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	opened, closed := state.Diff(prev, snap)
package state
