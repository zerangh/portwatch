package history

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
)

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func makeEntries() []Entry {
	return []Entry{
		{
			Timestamp: t0,
			Diff:      state.Diff{},
		},
		{
			Timestamp: t0.Add(time.Hour),
			Diff:      state.Diff{Added: []int{80, 443}, Removed: []int{22}},
		},
		{
			Timestamp: t0.Add(2 * time.Hour),
			Diff:      state.Diff{Added: []int{8080}},
		},
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize(nil)
	if s.Total != 0 {
		t.Errorf("expected 0 total, got %d", s.Total)
	}
	if s.ChangeRate() != 0 {
		t.Errorf("expected 0 change rate")
	}
}

func TestSummarize_Counts(t *testing.T) {
	s := Summarize(makeEntries())
	if s.Total != 3 {
		t.Errorf("expected 3, got %d", s.Total)
	}
	if s.WithChanges != 2 {
		t.Errorf("expected 2 with changes, got %d", s.WithChanges)
	}
}

func TestSummarize_MostAdded(t *testing.T) {
	s := Summarize(makeEntries())
	if s.MostAdded != 2 {
		t.Errorf("expected MostAdded=2, got %d", s.MostAdded)
	}
	if s.MostRemoved != 1 {
		t.Errorf("expected MostRemoved=1, got %d", s.MostRemoved)
	}
}

func TestSummarize_TimeRange(t *testing.T) {
	s := Summarize(makeEntries())
	if !s.FirstSeen.Equal(t0) {
		t.Errorf("unexpected FirstSeen: %v", s.FirstSeen)
	}
	if !s.LastSeen.Equal(t0.Add(2 * time.Hour)) {
		t.Errorf("unexpected LastSeen: %v", s.LastSeen)
	}
}

func TestChangeRate(t *testing.T) {
	s := Summarize(makeEntries())
	got := s.ChangeRate()
	want := 2.0 / 3.0
	if got < want-0.001 || got > want+0.001 {
		t.Errorf("expected ~%.4f, got %.4f", want, got)
	}
}
