package history

import (
	"testing"
	"time"
)

func makeEntry(ago time.Duration, ports []int) Entry {
	return Entry{
		Timestamp: time.Now().Add(-ago),
		Ports:     ports,
	}
}

func TestDefaultRetentionPolicy(t *testing.T) {
	p := DefaultRetentionPolicy()
	if p.MaxEntries != 100 {
		t.Errorf("expected MaxEntries=100, got %d", p.MaxEntries)
	}
	if p.MaxAge != 30*24*time.Hour {
		t.Errorf("unexpected MaxAge: %v", p.MaxAge)
	}
}

func TestApply_Empty(t *testing.T) {
	p := DefaultRetentionPolicy()
	result := p.Apply(nil)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestApply_MaxEntries(t *testing.T) {
	p := RetentionPolicy{MaxEntries: 3}
	entries := []Entry{
		makeEntry(1*time.Minute, []int{80}),
		makeEntry(2*time.Minute, []int{443}),
		makeEntry(3*time.Minute, []int{8080}),
		makeEntry(4*time.Minute, []int{22}),
		makeEntry(5*time.Minute, []int{3306}),
	}
	result := p.Apply(entries)
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
	// Newest should be first.
	if result[0].Ports[0] != 80 {
		t.Errorf("expected newest entry first, got port %d", result[0].Ports[0])
	}
}

func TestApply_MaxAge(t *testing.T) {
	p := RetentionPolicy{MaxAge: 10 * time.Minute}
	entries := []Entry{
		makeEntry(5*time.Minute, []int{80}),
		makeEntry(15*time.Minute, []int{443}),
		makeEntry(20*time.Minute, []int{22}),
	}
	result := p.Apply(entries)
	if len(result) != 1 {
		t.Errorf("expected 1 entry within age limit, got %d", len(result))
	}
	if result[0].Ports[0] != 80 {
		t.Errorf("expected port 80, got %d", result[0].Ports[0])
	}
}

func TestApply_AllPruned(t *testing.T) {
	p := RetentionPolicy{MaxAge: 1 * time.Minute}
	entries := []Entry{
		makeEntry(2*time.Minute, []int{80}),
		makeEntry(3*time.Minute, []int{443}),
	}
	result := p.Apply(entries)
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}
