package state

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot represents the recorded open ports at a point in time.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []int     `json:"ports"`
}

// Save writes a snapshot to the given file path as JSON.
func Save(path string, snap Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(snap)
}

// Load reads a snapshot from the given file path.
func Load(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return Snapshot{}, err
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// Diff computes ports opened and closed between two snapshots.
func Diff(prev, curr Snapshot) (opened, closed []int) {
	prevSet := toSet(prev.Ports)
	currSet := toSet(curr.Ports)

	for p := range currSet {
		if !prevSet[p] {
			opened = append(opened, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			closed = append(closed, p)
		}
	}
	return
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
