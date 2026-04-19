package labelmap_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/labelmap"
)

func writeLabelFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "labels.txt")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	lm, err := labelmap.Load("/nonexistent/labels.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := lm.Lookup(80); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestLoad_ParsesLabels(t *testing.T) {
	p := writeLabelFile(t, "# comment\n80 http\n443 https\n22 ssh\n")
	lm, err := labelmap.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cases := map[int]string{80: "http", 443: "https", 22: "ssh"}
	for port, want := range cases {
		if got := lm.Lookup(port); got != want {
			t.Errorf("port %d: got %q, want %q", port, got, want)
		}
	}
}

func TestLookup_NoMatch(t *testing.T) {
	lm := labelmap.New()
	if got := lm.Lookup(9999); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestSet_OverwritesLabel(t *testing.T) {
	lm := labelmap.New()
	lm.Set(80, "http")
	lm.Set(80, "web")
	if got := lm.Lookup(80); got != "web" {
		t.Errorf("got %q, want web", got)
	}
}

func TestAnnotate_ReturnsMatchingPorts(t *testing.T) {
	lm := labelmap.New()
	lm.Set(80, "http")
	lm.Set(443, "https")

	result := lm.Annotate([]int{80, 8080, 443})
	if result[80] != "http" {
		t.Errorf("expected http for 80")
	}
	if result[443] != "https" {
		t.Errorf("expected https for 443")
	}
	if _, ok := result[8080]; ok {
		t.Errorf("8080 should not be in result")
	}
}
