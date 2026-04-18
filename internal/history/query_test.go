package history

import (
	"testing"
	"time"
)

func baseTime() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func populatedHistory() *History {
	h := &History{}
	t0 := baseTime()
	h.entries = []Entry{
		{Timestamp: t0, OpenPorts: []int{80}},
		{Timestamp: t0.Add(time.Hour), OpenPorts: []int{80, 443}, Added: []int{443}},
		{Timestamp: t0.Add(2 * time.Hour), OpenPorts: []int{80}, Removed: []int{443}},
		{Timestamp: t0.Add(3 * time.Hour), OpenPorts: []int{80}},
	}
	return h
}

func TestFilter_Since(t *testing.T) {
	h := populatedHistory()
	res := h.Since(baseTime().Add(90 * time.Minute))
	if len(res) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res))
	}
}

func TestFilter_Limit(t *testing.T) {
	h := populatedHistory()
	res := h.Filter(Query{Limit: 2})
	if len(res) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res))
	}
}

func TestFilter_HasDiff(t *testing.T) {
	h := populatedHistory()
	res := h.WithChanges()
	if len(res) != 2 {
		t.Fatalf("expected 2 entries with changes, got %d", len(res))
	}
	for _, e := range res {
		if len(e.Added) == 0 && len(e.Removed) == 0 {
			t.Error("entry without changes returned by WithChanges")
		}
	}
}

func TestFilter_Until(t *testing.T) {
	h := populatedHistory()
	res := h.Filter(Query{Until: baseTime().Add(time.Hour)})
	if len(res) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res))
	}
}
