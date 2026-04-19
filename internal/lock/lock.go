// Package lock provides a simple file-based lock to prevent concurrent
// portwatch processes from running against the same state file.
package lock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ErrLocked is returned when a lock file already exists and the owning
// process is still running.
var ErrLocked = errors.New("lock: already locked by another process")

// Lock represents a file-based process lock.
type Lock struct {
	path string
}

// New returns a Lock whose lock file is derived from statePath.
func New(statePath string) *Lock {
	dir := filepath.Dir(statePath)
	base := filepath.Base(statePath)
	return &Lock{path: filepath.Join(dir, "."+base+".lock")}
}

// Acquire creates the lock file containing the current PID.
// Returns ErrLocked if a live process already holds the lock.
func (l *Lock) Acquire() error {
	if data, err := os.ReadFile(l.path); err == nil {
		pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err == nil && processAlive(pid) {
			return fmt.Errorf("%w (pid %d, lock %s)", ErrLocked, pid, l.path)
		}
		// Stale lock — remove it.
		_ = os.Remove(l.path)
	}
	return os.WriteFile(l.path, []byte(strconv.Itoa(os.Getpid())), 0o600)
}

// Release removes the lock file.
func (l *Lock) Release() error {
	return os.Remove(l.path)
}

// Path returns the path of the lock file.
func (l *Lock) Path() string { return l.path }

// processAlive reports whether a process with the given PID exists.
func processAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Unix, FindProcess always succeeds; signal 0 checks existence.
	return proc.Signal(os.Signal(nil)) == nil
}
