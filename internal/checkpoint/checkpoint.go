// Package checkpoint persists the last-known scan time so portwatch
// can detect whether a scan was skipped across restarts.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Checkpoint holds the metadata written to disk after each scan.
type Checkpoint struct {
	LastScan  time.Time `json:"last_scan"`
	PortCount int       `json:"port_count"`
}

// Manager reads and writes checkpoint files.
type Manager struct {
	path string
}

// New returns a Manager that stores checkpoints at path.
func New(path string) (*Manager, error) {
	if path == "" {
		return nil, errors.New("checkpoint: path must not be empty")
	}
	return &Manager{path: path}, nil
}

// Save writes cp to disk, overwriting any existing checkpoint.
func (m *Manager) Save(cp Checkpoint) error {
	f, err := os.CreateTemp("", "checkpoint-*")
	if err != nil {
		return err
	}
	if err := json.NewEncoder(f).Encode(cp); err != nil {
		f.Close()
		return err
	}
	f.Close()
	return os.Rename(f.Name(), m.path)
}

// Load reads the checkpoint from disk. If the file does not exist a
// zero-value Checkpoint and nil error are returned.
func (m *Manager) Load() (Checkpoint, error) {
	var cp Checkpoint
	f, err := os.Open(m.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cp, nil
		}
		return cp, err
	}
	defer f.Close()
	return cp, json.NewDecoder(f).Decode(&cp)
}

// Age returns the duration since the last saved scan.
// If no checkpoint exists, Age returns zero and ok=false.
func (m *Manager) Age() (time.Duration, bool) {
	cp, err := m.Load()
	if err != nil || cp.LastScan.IsZero() {
		return 0, false
	}
	return time.Since(cp.LastScan), true
}
