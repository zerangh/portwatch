package portrank

import (
	"testing"
)

func TestRank_Builtin(t *testing.T) {
	r := New(nil)
	if got := r.Rank(22); got != RankCritical {
		t.Fatalf("expected critical for port 22, got %s", got)
	}
	if got := r.Rank(80); got != RankHigh {
		t.Fatalf("expected high for port 80, got %s", got)
	}
	if got := r.Rank(9999); got != RankLow {
		t.Fatalf("expected low for unknown port, got %s", got)
	}
}

func TestRank_Override(t *testing.T) {
	r := New(map[int]Rank{9999: RankCritical})
	if got := r.Rank(9999); got != RankCritical {
		t.Fatalf("expected override critical, got %s", got)
	}
	// builtin still works
	if got := r.Rank(22); got != RankCritical {
		t.Fatalf("expected critical for 22, got %s", got)
	}
}

func TestRankString(t *testing.T) {
	cases := []struct {
		rank Rank
		want string
	}{
		{RankLow, "low"},
		{RankMedium, "medium"},
		{RankHigh, "high"},
		{RankCritical, "critical"},
		{Rank(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.rank.String(); got != c.want {
			t.Errorf("Rank(%d).String() = %q, want %q", c.rank, got, c.want)
		}
	}
}

func TestSortByRank(t *testing.T) {
	r := New(nil)
	ports := []int{9999, 22, 80, 6379}
	sorted := r.SortByRank(ports)
	if sorted[0] != 22 {
		t.Fatalf("expected 22 first (critical), got %d", sorted[0])
	}
	if sorted[len(sorted)-1] != 9999 {
		t.Fatalf("expected 9999 last (low), got %d", sorted[len(sorted)-1])
	}
}

func TestFilterByMinRank(t *testing.T) {
	r := New(nil)
	ports := []int{22, 80, 6379, 9999}
	got := r.FilterByMinRank(ports, RankHigh)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports at or above high, got %d: %v", len(got), got)
	}
}
