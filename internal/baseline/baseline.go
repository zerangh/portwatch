// Package baseline manages a known-good snapshot of open ports that can be
// compared against current scan results to suppress expected ports from alerts.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Baseline holds a saved snapshot of open ports.
type Baseline struct {
	Ports     []int     `json:"ports"`
	CreatedAt time.Time `json:"created_at"`
	Path      string    `json:"-"`
}

// New returns an empty Baseline backed by the given file path.
func New(path string) *Baseline {
	return &Baseline{Path: path}
}

// Save writes the baseline to disk.
func (b *Baseline) Save(ports []int) error {
	if b.Path == "" {
		return errors.New("baseline: path is empty")
	}
	b.Ports = ports
	b.CreatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.Path, data, 0o644)
}

// Load reads the baseline from disk. Returns an empty Baseline if the file
// does not exist.
func (b *Baseline) Load() error {
	data, err := os.ReadFile(b.Path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, b)
}

// Filter removes ports that are present in the baseline from the provided
// slice, returning only unexpected ports.
func (b *Baseline) Filter(ports []int) []int {
	known := make(map[int]struct{}, len(b.Ports))
	for _, p := range b.Ports {
		known[p] = struct{}{}
	}
	out := ports[:0:0]
	for _, p := range ports {
		if _, ok := known[p]; !ok {
			out = append(out, p)
		}
	}
	return out
}

// IsEmpty reports whether the baseline contains no ports.
func (b *Baseline) IsEmpty() bool {
	return len(b.Ports) == 0
}
