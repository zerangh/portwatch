package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
)

func TestApply_NoRules(t *testing.T) {
	f := filter.New(filter.Rule{})
	ports := []int{80, 443, 8080}
	got := f.Apply(ports)
	if len(got) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(got))
	}
}

func TestApply_ExcludePorts(t *testing.T) {
	f := filter.New(filter.Rule{ExcludePorts: []int{80, 443}})
	ports := []int{80, 443, 8080}
	got := f.Apply(ports)
	if len(got) != 1 || got[0] != 8080 {
		t.Fatalf("expected [8080], got %v", got)
	}
}

func TestApply_ExcludeRange(t *testing.T) {
	f := filter.New(filter.Rule{
		ExcludeRanges: []filter.Range{{Low: 1000, High: 2000}},
	})
	ports := []int{80, 1500, 3000}
	got := f.Apply(ports)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %v", got)
	}
	for _, p := range got {
		if p == 1500 {
			t.Fatal("port 1500 should have been excluded")
		}
	}
}

func TestApply_EmptyPorts(t *testing.T) {
	f := filter.New(filter.Rule{ExcludePorts: []int{80}})
	got := f.Apply([]int{})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}

func TestApply_CombinedRules(t *testing.T) {
	f := filter.New(filter.Rule{
		ExcludePorts:  []int{22},
		ExcludeRanges: []filter.Range{{Low: 8000, High: 9000}},
	})
	ports := []int{22, 80, 8080, 9090}
	got := f.Apply(ports)
	if len(got) != 2 {
		t.Fatalf("expected [80 9090], got %v", got)
	}
}
