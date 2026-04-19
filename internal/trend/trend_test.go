package trend_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/trend"
)

func makeEntry(added, removed []int, ago time.Duration) history.Entry {
	return history.Entry{
		Timestamp: time.Now().Add(-ago),
		Diff: state.Diff{
			Added:   added,
			Removed: removed,
		},
	}
}

func TestAnalyze_Empty(t *testing.T) {
	r := trend.Analyze(nil, time.Hour)
	if r.Direction != trend.Stable {
		t.Errorf("expected Stable, got %s", r.Direction)
	}
	if r.Samples != 0 {
		t.Errorf("expected 0 samples")
	}
}

func TestAnalyze_Growing(t *testing.T) {
	entries := []history.Entry{
		makeEntry([]int{80, 443}, nil, 10*time.Minute),
		makeEntry([]int{8080}, nil, 5*time.Minute),
	}
	r := trend.Analyze(entries, time.Hour)
	if r.Direction != trend.Growing {
		t.Errorf("expected Growing, got %s", r.Direction)
	}
	if r.NetChange != 3 {
		t.Errorf("expected net 3, got %d", r.NetChange)
	}
}

func TestAnalyze_Shrinking(t *testing.T) {
	entries := []history.Entry{
		makeEntry(nil, []int{80, 443, 8080}, 5*time.Minute),
	}
	r := trend.Analyze(entries, time.Hour)
	if r.Direction != trend.Shrinking {
		t.Errorf("expected Shrinking, got %s", r.Direction)
	}
}

func TestAnalyze_ExcludesOldEntries(t *testing.T) {
	entries := []history.Entry{
		makeEntry([]int{80}, nil, 2*time.Hour), // outside window
		makeEntry(nil, []int{443}, 10*time.Minute),
	}
	r := trend.Analyze(entries, time.Hour)
	if r.Samples != 1 {
		t.Errorf("expected 1 sample, got %d", r.Samples)
	}
	if r.Direction != trend.Shrinking {
		t.Errorf("expected Shrinking")
	}
}

func TestAnalyze_Stable(t *testing.T) {
	entries := []history.Entry{
		makeEntry([]int{80}, []int{80}, 5*time.Minute),
	}
	r := trend.Analyze(entries, time.Hour)
	if r.Direction != trend.Stable {
		t.Errorf("expected Stable, got %s", r.Direction)
	}
}
