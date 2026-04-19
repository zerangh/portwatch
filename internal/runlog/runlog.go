// Package runlog records a rolling log of scan run outcomes.
package runlog

import (
	"encoding/json"
	"os"
	"time"
)

// Entry represents a single scan run outcome.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	PortsFound int       `json:"ports_found"`
	Changed    bool      `json:"changed"`
	Error      string    `json:"error,omitempty"`
	DurationMs int64     `json:"duration_ms"`
}

// RunLog holds a capped list of run entries persisted to disk.
type RunLog struct {
	path    string
	maxSize int
}

// New returns a RunLog writing to path, keeping at most maxSize entries.
func New(path string, maxSize int) (*RunLog, error) {
	if path == "" {
		return nil, fmt.Errorf("runlog: path must not be empty")
	}
	if maxSize <= 0 {
		maxSize = 100
	}
	return &RunLog{path: path, maxSize: maxSize}, nil
}

// Append adds an entry to the log, pruning old entries if needed.
func (r *RunLog) Append(e Entry) error {
	entries, _ := r.Load()
	entries = append(entries, e)
	if len(entries) > r.maxSize {
		entries = entries[len(entries)-r.maxSize:]
	}
	return r.save(entries)
}

// Load returns all stored entries.
func (r *RunLog) Load() ([]Entry, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *RunLog) save(entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0o644)
}
