// Package lock provides file-based mutual exclusion for portwatch processes.
//
// A Lock is tied to a state file path; the lock file is placed in the same
// directory with a dot-prefixed name so that it is hidden on Unix systems.
//
// Usage:
//
//	l := lock.New(cfg.StatePath)
//	if err := l.Acquire(); err != nil {
//		log.Fatal(err)
//	}
//	defer l.Release()
//
// Stale locks (whose owning process is no longer alive) are automatically
// removed and re-acquired on the next Acquire call.
package lock
