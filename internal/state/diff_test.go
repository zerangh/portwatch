package state

import (
	"strings"
	"testing"
)

func TestDiff_IsEmpty(t *testing.T) {
	d := Diff{}
	if !d.IsEmpty() {
		t.Error("expected empty diff to be empty")
	}
	d.Added = []int{80}
	if d.IsEmpty() {
		t.Error("expected non-empty diff")
	}
}

func TestDiff_Summary_NoChanges(t *testing.T) {
	d := Diff{}
	if d.Summary() != "no changes" {
		t.Errorf("unexpected summary: %s", d.Summary())
	}
}

func TestDiff_Summary_Added(t *testing.T) {
	d := Diff{Added: []int{8080, 443}}
	s := d.Summary()
	if !strings.Contains(s, "opened") {
		t.Errorf("expected 'opened' in summary: %s", s)
	}
	if !strings.Contains(s, "443") || !strings.Contains(s, "8080") {
		t.Errorf("expected ports in summary: %s", s)
	}
}

func TestDiff_Summary_Removed(t *testing.T) {
	d := Diff{Removed: []int{22}}
	s := d.Summary()
	if !strings.Contains(s, "closed") {
		t.Errorf("expected 'closed' in summary: %s", s)
	}
}

func TestDiff_Summary_Both(t *testing.T) {
	d := Diff{Added: []int{9000}, Removed: []int{80}}
	s := d.Summary()
	if !strings.Contains(s, "opened") || !strings.Contains(s, "closed") {
		t.Errorf("expected both in summary: %s", s)
	}
}
