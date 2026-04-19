package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/checkpoint"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_EmptyPathReturnsError(t *testing.T) {
	_, err := checkpoint.New("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLoad_MissingFile_ReturnsZero(t *testing.T) {
	m, _ := checkpoint.New(tmpPath(t))
	cp, err := m.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cp.LastScan.IsZero() {
		t.Error("expected zero LastScan for missing file")
	}
}

func TestSaveAndLoad(t *testing.T) {
	m, _ := checkpoint.New(tmpPath(t))
	now := time.Now().Truncate(time.Second)
	want := checkpoint.Checkpoint{LastScan: now, PortCount: 42}

	if err := m.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := m.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !got.LastScan.Equal(want.LastScan) {
		t.Errorf("LastScan: got %v, want %v", got.LastScan, want.LastScan)
	}
	if got.PortCount != want.PortCount {
		t.Errorf("PortCount: got %d, want %d", got.PortCount, want.PortCount)
	}
}

func TestAge_NoCheckpoint(t *testing.T) {
	m, _ := checkpoint.New(tmpPath(t))
	_, ok := m.Age()
	if ok {
		t.Error("expected ok=false when no checkpoint exists")
	}
}

func TestAge_ReturnsElapsed(t *testing.T) {
	m, _ := checkpoint.New(tmpPath(t))
	past := time.Now().Add(-5 * time.Second)
	_ = m.Save(checkpoint.Checkpoint{LastScan: past, PortCount: 1})

	age, ok := m.Age()
	if !ok {
		t.Fatal("expected ok=true")
	}
	if age < 4*time.Second {
		t.Errorf("age too small: %v", age)
	}
}

func TestSave_Atomic(t *testing.T) {
	path := tmpPath(t)
	m, _ := checkpoint.New(path)
	_ = m.Save(checkpoint.Checkpoint{PortCount: 1})
	_ = m.Save(checkpoint.Checkpoint{PortCount: 2})

	cp, _ := m.Load()
	if cp.PortCount != 2 {
		t.Errorf("expected PortCount=2, got %d", cp.PortCount)
	}
	// no temp files left behind
	entries, _ := os.ReadDir(filepath.Dir(path))
	if len(entries) != 1 {
		t.Errorf("expected 1 file, found %d", len(entries))
	}
}
