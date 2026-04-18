// Package snapshot provides types and utilities for capturing and comparing
// the state of open ports on a host at a specific point in time.
//
// A Snapshot records the host name, the list of open ports (sorted), and the
// time at which the capture occurred.
//
// Use New to create a snapshot, Contains to check for a specific port, and
// Compare to determine which ports were added or removed relative to a
// previous snapshot.
//
// Example:
//
//	curr := snapshot.New("localhost", openPorts)
//	added, removed := curr.Compare(prev)
package snapshot
