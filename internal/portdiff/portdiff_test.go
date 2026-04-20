package portdiff_test

import (
	"testing"

	"github.com/user/portwatch/internal/portdiff"
)

func TestCompute_NoChange(t *testing.T) {
	r := portdiff.Compute([]int{80, 443}, []int{80, 443})
	if !r.IsEmpty() {
		t.Fatalf("expected empty diff, got added=%v removed=%v", r.Added, r.Removed)
	}
}

func TestCompute_Added(t *testing.T) {
	r := portdiff.Compute([]int{80}, []int{80, 8080})
	if len(r.Added) != 1 || r.Added[0] != 8080 {
		t.Fatalf("expected [8080] added, got %v", r.Added)
	}
	if len(r.Removed) != 0 {
		t.Fatalf("expected no removed, got %v", r.Removed)
	}
}

func TestCompute_Removed(t *testing.T) {
	r := portdiff.Compute([]int{80, 443}, []int{80})
	if len(r.Removed) != 1 || r.Removed[0] != 443 {
		t.Fatalf("expected [443] removed, got %v", r.Removed)
	}
	if len(r.Added) != 0 {
		t.Fatalf("expected no added, got %v", r.Added)
	}
}

func TestCompute_Mixed(t *testing.T) {
	r := portdiff.Compute([]int{22, 80}, []int{80, 443})
	if len(r.Added) != 1 || r.Added[0] != 443 {
		t.Fatalf("expected [443] added, got %v", r.Added)
	}
	if len(r.Removed) != 1 || r.Removed[0] != 22 {
		t.Fatalf("expected [22] removed, got %v", r.Removed)
	}
}

func TestCompute_Empty(t *testing.T) {
	r := portdiff.Compute(nil, nil)
	if !r.IsEmpty() {
		t.Fatal("expected empty diff for nil inputs")
	}
}

func TestSummary_NoChanges(t *testing.T) {
	r := portdiff.Result{}
	if r.Summary() != "no changes" {
		t.Fatalf("unexpected summary: %q", r.Summary())
	}
}

func TestSummary_Added(t *testing.T) {
	r := portdiff.Result{Added: []int{80, 443}}
	if r.Summary() != "+2 added" {
		t.Fatalf("unexpected summary: %q", r.Summary())
	}
}

func TestSummary_Both(t *testing.T) {
	r := portdiff.Result{Added: []int{8080}, Removed: []int{22, 23}}
	want := "+1 added, -2 removed"
	if r.Summary() != want {
		t.Fatalf("expected %q, got %q", want, r.Summary())
	}
}
