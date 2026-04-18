package state

import (
	"os"
	"sort"
	"testing"
	"time"
)

func TestSaveAndLoad(t *testing.T) {
	tmp, err := os.CreateTemp("", "portwatch-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	snap := Snapshot{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Ports:     []int{80, 443, 8080},
	}

	if err := Save(tmp.Name(), snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(tmp.Name())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if !loaded.Timestamp.Equal(snap.Timestamp) {
		t.Errorf("timestamp mismatch: got %v want %v", loaded.Timestamp, snap.Timestamp)
	}
	if len(loaded.Ports) != len(snap.Ports) {
		t.Errorf("ports length mismatch: got %d want %d", len(loaded.Ports), len(snap.Ports))
	}
}

func TestDiff(t *testing.T) {
	prev := Snapshot{Ports: []int{80, 443, 22}}
	curr := Snapshot{Ports: []int{80, 8080, 22}}

	opened, closed := Diff(prev, curr)

	sort.Ints(opened)
	sort.Ints(closed)

	if len(opened) != 1 || opened[0] != 8080 {
		t.Errorf("opened: expected [8080], got %v", opened)
	}
	if len(closed) != 1 || closed[0] != 443 {
		t.Errorf("closed: expected [443], got %v", closed)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/portwatch.json")
	if err == nil {
		t.Error("expected error loading missing file, got nil")
	}
}
