package runlog_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/runlog"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "runlog.json")
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	rl, _ := runlog.New(tmpPath(t), 10)
	entries, err := rl.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty, got %d", len(entries))
	}
}

func TestAppendAndLoad(t *testing.T) {
	rl, _ := runlog.New(tmpPath(t), 10)
	e := runlog.Entry{
		Timestamp:  time.Now().UTC(),
		PortsFound: 3,
		Changed:    true,
		DurationMs: 42,
	}
	if err := rl.Append(e); err != nil {
		t.Fatalf("append: %v", err)
	}
	entries, err := rl.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].PortsFound != 3 {
		t.Errorf("ports found mismatch")
	}
}

func TestAppend_CapsAtMaxSize(t *testing.T) {
	rl, _ := runlog.New(tmpPath(t), 5)
	for i := 0; i < 8; i++ {
		_ = rl.Append(runlog.Entry{Timestamp: time.Now().UTC(), PortsFound: i})
	}
	entries, _ := rl.Load()
	if len(entries) != 5 {
		t.Fatalf("expected 5, got %d", len(entries))
	}
	if entries[0].PortsFound != 3 {
		t.Errorf("expected oldest kept entry to have PortsFound=3")
	}
}

func TestNew_EmptyPathReturnsError(t *testing.T) {
	_, err := runlog.New("", 10)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestAppend_ErrorEntry(t *testing.T) {
	rl, _ := runlog.New(tmpPath(t), 10)
	_ = rl.Append(runlog.Entry{Timestamp: time.Now().UTC(), Error: "scan failed"})
	entries, _ := rl.Load()
	if entries[0].Error != "scan failed" {
		t.Errorf("error field not persisted")
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	p := tmpPath(t)
	_ = os.WriteFile(p, []byte("not json"), 0o644)
	rl, _ := runlog.New(p, 10)
	_, err := rl.Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
