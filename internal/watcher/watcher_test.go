package watcher_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watcher"
)

func newTestScanner(t *testing.T, start, end int) *scanner.Scanner {
	t.Helper()
	s, err := scanner.New(scanner.Config{
		Host:        "127.0.0.1",
		PortStart:   start,
		PortEnd:     end,
		Concurrency: 10,
	})
	if err != nil {
		t.Fatalf("scanner.New: %v", err)
	}
	return s
}

func TestWatcher_RunCreatesStateFile(t *testing.T) {
	dir := t.TempDir()
	statePath := filepath.Join(dir, "state.json")

	w, err := watcher.New(watcher.Config{
		Scanner:   newTestScanner(t, 65400, 65410),
		StatePath: statePath,
		Writer:    &bytes.Buffer{},
		Format:    "text",
	})
	if err != nil {
		t.Fatalf("watcher.New: %v", err)
	}

	if err := w.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Error("expected state file to exist after Run")
	}
}

func TestNew_NilScannerReturnsError(t *testing.T) {
	_, err := watcher.New(watcher.Config{
		Scanner:   nil,
		StatePath: "/tmp/state.json",
	})
	if err == nil {
		t.Error("expected error for nil scanner")
	}
}

func TestNew_EmptyStatePathReturnsError(t *testing.T) {
	_, err := watcher.New(watcher.Config{
		Scanner:   newTestScanner(t, 65400, 65410),
		StatePath: "",
	})
	if err == nil {
		t.Error("expected error for empty state path")
	}
}
