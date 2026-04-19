package suppress

import (
	"os"
	"path/filepath"
	"testing"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "suppress.json")
}

func TestNew_EmptyPathReturnsError(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNew_MissingFile_ReturnsEmpty(t *testing.T) {
	l, err := New(tmpPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Contains(80) {
		t.Fatal("expected empty list")
	}
}

func TestAddAndContains(t *testing.T) {
	l, _ := New(tmpPath(t))
	if err := l.Add(443); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if !l.Contains(443) {
		t.Fatal("expected 443 to be suppressed")
	}
	if l.Contains(80) {
		t.Fatal("80 should not be suppressed")
	}
}

func TestRemove(t *testing.T) {
	l, _ := New(tmpPath(t))
	_ = l.Add(8080)
	if err := l.Remove(8080); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if l.Contains(8080) {
		t.Fatal("expected 8080 to be removed")
	}
}

func TestFilter(t *testing.T) {
	l, _ := New(tmpPath(t))
	_ = l.Add(22)
	_ = l.Add(443)
	got := l.Filter([]int{22, 80, 443, 8080})
	if len(got) != 2 || got[0] != 80 || got[1] != 8080 {
		t.Fatalf("unexpected filter result: %v", got)
	}
}

func TestPersistence(t *testing.T) {
	p := tmpPath(t)
	l1, _ := New(p)
	_ = l1.Add(9000)

	l2, err := New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !l2.Contains(9000) {
		t.Fatal("expected 9000 to persist across reload")
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	p := tmpPath(t)
	_ = os.WriteFile(p, []byte("not json"), 0o644)
	_, err := New(p)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
