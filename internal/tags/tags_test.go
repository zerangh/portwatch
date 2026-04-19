package tags_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/tags"
)

func writeTags(t *testing.T, data []tags.Tag) string {
	t.Helper()
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "tags.json")
	if err := os.WriteFile(p, b, 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLookup_ExactPort(t *testing.T) {
	m := tags.New()
	m.Add(tags.Tag{Label: "http", Ports: []int{80, 8080}})
	labels := m.Lookup(80)
	if len(labels) != 1 || labels[0] != "http" {
		t.Fatalf("expected [http], got %v", labels)
	}
}

func TestLookup_Range(t *testing.T) {
	m := tags.New()
	m.Add(tags.Tag{Label: "ephemeral", From: 49152, To: 65535})
	if labels := m.Lookup(50000); len(labels) == 0 {
		t.Fatal("expected match for port in range")
	}
	if labels := m.Lookup(80); len(labels) != 0 {
		t.Fatalf("unexpected match: %v", labels)
	}
}

func TestLookup_NoMatch(t *testing.T) {
	m := tags.New()
	m.Add(tags.Tag{Label: "ssh", Ports: []int{22}})
	if labels := m.Lookup(443); len(labels) != 0 {
		t.Fatalf("expected no match, got %v", labels)
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	m, err := tags.Load("/nonexistent/tags.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Len() != 0 {
		t.Fatalf("expected empty map, got len %d", m.Len())
	}
}

func TestLoad_ValidFile(t *testing.T) {
	path := writeTags(t, []tags.Tag{
		{Label: "dns", Ports: []int{53}},
		{Label: "web", From: 80, To: 81},
	})
	m, err := tags.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Len() != 2 {
		t.Fatalf("expected 2 tags, got %d", m.Len())
	}
	if labels := m.Lookup(53); len(labels) == 0 || labels[0] != "dns" {
		t.Fatalf("expected dns label, got %v", labels)
	}
}

func TestLookup_MultipleLabels(t *testing.T) {
	m := tags.New()
	m.Add(tags.Tag{Label: "web", Ports: []int{443}})
	m.Add(tags.Tag{Label: "tls", From: 400, To: 500})
	labels := m.Lookup(443)
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %v", labels)
	}
}
