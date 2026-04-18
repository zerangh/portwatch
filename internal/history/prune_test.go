package history

import (
	"testing"
	"time"
)

func TestPrune_RemovesOldEntries(t *testing.T) {
	path := tmpPath(t)
	h, _ := New(path)

	now := time.Now()
	old := Entry{Timestamp: now.Add(-48 * time.Hour), OpenPorts: []int{80}}
	recent := Entry{Timestamp: now.Add(-1 * time.Hour), OpenPorts: []int{443}}

	_ = h.Append(old)
	_ = h.Append(recent)

	policy := DefaultRetentionPolicy()
	policy.MaxAge = 24 * time.Hour

	removed, err := Prune(path, policy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}

	entries, _ := h.Load()
	if len(entries) != 1 {
		t.Errorf("expected 1 remaining entry, got %d", len(entries))
	}
	if len(entries) > 0 && entries[0].OpenPorts[0] != 443 {
		t.Errorf("expected remaining entry to have port 443")
	}
}

func TestPrune_NothingToRemove(t *testing.T) {
	path := tmpPath(t)
	h, _ := New(path)

	_ = h.Append(Entry{Timestamp: time.Now(), OpenPorts: []int{80}})

	policy := DefaultRetentionPolicy()
	policy.MaxAge = 24 * time.Hour

	removed, err := Prune(path, policy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 0 {
		t.Errorf("expected 0 removed, got %d", removed)
	}
}

func TestPrune_MissingFile(t *testing.T) {
	_, err := Prune("/nonexistent/path/history.json", DefaultRetentionPolicy())
	if err == nil {
		t.Error("expected error for missing file")
	}
}
