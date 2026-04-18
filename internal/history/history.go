// Package history provides persistent scan history tracking for portwatch.
package history

import (
	"encoding/json"
	"os"
	"time"
)

// Entry represents a single scan result recorded in history.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	OpenPorts []int     `json:"open_ports"`
	Added     []int     `json:"added,omitempty"`
	Removed   []int     `json:"removed,omitempty"`
}

// History holds a list of scan entries.
type History struct {
	Entries []Entry `json:"entries"`
	path    string
}

// New creates a new History backed by the given file path.
func New(path string) *History {
	return &History{path: path}
}

// Load reads history from disk. Returns empty history if file does not exist.
func (h *History) Load() error {
	data, err := os.ReadFile(h.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, h)
}

// Append adds a new entry and persists history to disk.
func (h *History) Append(e Entry) error {
	h.Entries = append(h.Entries, e)
	return h.save()
}

// Last returns the most recent entry, or nil if history is empty.
func (h *History) Last() *Entry {
	if len(h.Entries) == 0 {
		return nil
	}
	e := h.Entries[len(h.Entries)-1]
	return &e
}

func (h *History) save() error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o644)
}
