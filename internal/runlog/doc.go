// Package runlog provides a persistent, capped log of portwatch scan run
// outcomes. Each entry records the time, number of open ports found, whether
// a change was detected, any error encountered, and the scan duration.
//
// Entries are stored as a JSON array on disk. When the log grows beyond the
// configured maximum size the oldest entries are automatically pruned so the
// file stays bounded.
//
// Typical usage:
//
//	rl, err := runlog.New("/var/lib/portwatch/runs.json", 200)
//	if err != nil { ... }
//	_ = rl.Append(runlog.Entry{
//		Timestamp:  time.Now().UTC(),
//		PortsFound: 12,
//		Changed:    true,
//		DurationMs: 310,
//	})
package runlog
