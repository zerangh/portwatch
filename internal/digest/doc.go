// Package digest provides lightweight fingerprinting of port sets.
//
// A Digest is a short hex string derived from a sorted, comma-joined
// list of port numbers hashed with SHA-256. It allows the watcher and
// reporter to quickly detect whether a scan result has changed since
// the last run without performing a full state diff.
//
// Usage:
//
//	d := digest.FromPorts([]int{22, 80, 443})
//	if !digest.Equal(prev, d) {
//	    // ports have changed, run full diff
//	}
package digest
