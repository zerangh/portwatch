package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/baseline"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestSaveAndLoad(t *testing.T) {
	p := tmpPath(t)
	b := baseline.New(p)
	ports := []int{22, 80, 443}
	if err := b.Save(ports); err != nil {
		t.Fatalf("Save: %v", err)
	}
	b2 := baseline.New(p)
	if err := b2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(b2.Ports) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(b2.Ports))
	}
	if b2.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if b2.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt is in the future")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	b := baseline.New(tmpPath(t))
	if err := b.Load(); err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if !b.IsEmpty() {
		t.Error("expected empty baseline")
	}
}

func TestFilter_RemovesKnownPorts(t *testing.T) {
	b := baseline.New("")
	b.Ports = []int{22, 80}
	result := b.Filter([]int{22, 80, 8080, 9090})
	if len(result) != 2 {
		t.Fatalf("expected 2 ports, got %d: %v", len(result), result)
	}
	for _, p := range result {
		if p == 22 || p == 80 {
			t.Errorf("baseline port %d should have been filtered", p)
		}
	}
}

func TestFilter_EmptyBaseline(t *testing.T) {
	b := baseline.New("")
	ports := []int{22, 443}
	result := b.Filter(ports)
	if len(result) != len(ports) {
		t.Fatalf("expected all ports to pass through, got %v", result)
	}
}

func TestSave_EmptyPath(t *testing.T) {
	b := baseline.New("")
	if err := b.Save([]int{80}); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestIsEmpty(t *testing.T) {
	b := baseline.New("")
	if !b.IsEmpty() {
		t.Error("new baseline should be empty")
	}
	b.Ports = []int{80}
	if b.IsEmpty() {
		t.Error("baseline with ports should not be empty")
	}
}

func TestSaveCreatesFile(t *testing.T) {
	p := tmpPath(t)
	b := baseline.New(p)
	_ = b.Save([]int{8080})
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}
