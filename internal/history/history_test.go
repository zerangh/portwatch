package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	h := history.New(tmpPath(t))
	if err := h.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(h.Entries))
	}
}

func TestAppendAndLoad(t *testing.T) {
	path := tmpPath(t)
	h := history.New(path)

	e := history.Entry{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		OpenPorts: []int{80, 443},
		Added:     []int{443},
	}
	if err := h.Append(e); err != nil {
		t.Fatalf("append failed: %v", err)
	}

	h2 := history.New(path)
	if err := h2.Load(); err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(h2.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h2.Entries))
	}
	if h2.Entries[0].OpenPorts[1] != 443 {
		t.Errorf("unexpected port value")
	}
}

func TestLast_EmptyHistory(t *testing.T) {
	h := history.New(tmpPath(t))
	if h.Last() != nil {
		t.Error("expected nil for empty history")
	}
}

func TestLast_ReturnsNewest(t *testing.T) {
	path := tmpPath(t)
	h := history.New(path)
	_ = h.Append(history.Entry{OpenPorts: []int{22}})
	_ = h.Append(history.Entry{OpenPorts: []int{22, 80}})

	last := h.Last()
	if last == nil || len(last.OpenPorts) != 2 {
		t.Errorf("unexpected last entry: %+v", last)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := tmpPath(t)
	_ = os.WriteFile(path, []byte("not json{"), 0o644)
	h := history.New(path)
	if err := h.Load(); err == nil {
		t.Error("expected error for invalid JSON")
	}
}
