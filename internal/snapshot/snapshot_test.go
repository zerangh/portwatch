package snapshot

import (
	"testing"
	"time"
)

func TestNew_SortsPorts(t *testing.T) {
	s := New("localhost", []int{443, 80, 22})
	expected := []int{22, 80, 443}
	for i, p := range s.Ports {
		if p != expected[i] {
			t.Errorf("port[%d] = %d, want %d", i, p, expected[i])
		}
	}
}

func TestNew_SetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	s := New("localhost", []int{80})
	after := time.Now().UTC()
	if s.Timestamp.Before(before) || s.Timestamp.After(after) {
		t.Errorf("unexpected timestamp: %v", s.Timestamp)
	}
}

func TestNew_EmptyPorts(t *testing.T) {
	s := New("localhost", []int{})
	if len(s.Ports) != 0 {
		t.Errorf("expected empty ports, got %v", s.Ports)
	}
}

func TestContains(t *testing.T) {
	s := New("localhost", []int{22, 80, 443})
	if !s.Contains(80) {
		t.Error("expected Contains(80) = true")
	}
	if s.Contains(8080) {
		t.Error("expected Contains(8080) = false")
	}
}

func TestCompare_NilPrev(t *testing.T) {
	s := New("localhost", []int{22, 80})
	added, removed := s.Compare(nil)
	if len(added) != 2 {
		t.Errorf("added = %v, want 2 ports", added)
	}
	if len(removed) != 0 {
		t.Errorf("removed = %v, want empty", removed)
	}
}

func TestCompare_DetectsChanges(t *testing.T) {
	prev := New("localhost", []int{22, 80, 443})
	curr := New("localhost", []int{22, 443, 8080})
	added, removed := curr.Compare(prev)

	if len(added) != 1 || added[0] != 8080 {
		t.Errorf("added = %v, want [8080]", added)
	}
	if len(removed) != 1 || removed[0] != 80 {
		t.Errorf("removed = %v, want [80]", removed)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	prev := New("localhost", []int{22, 80})
	curr := New("localhost", []int{22, 80})
	added, removed := curr.Compare(prev)
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no changes, got added=%v removed=%v", added, removed)
	}
}

func TestString(t *testing.T) {
	s := New("myhost", []int{80})
	str := s.String()
	if str == "" {
		t.Error("expected non-empty string")
	}
}
