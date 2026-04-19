package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
)

func TestChain_Empty(t *testing.T) {
	c := filter.NewChain()
	ports := []int{80, 443}
	got := c.Apply(ports)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %v", got)
	}
}

func TestChain_SingleFilter(t *testing.T) {
	f := filter.New(filter.Rule{ExcludePorts: []int{80}})
	c := filter.NewChain(f)
	got := c.Apply([]int{80, 443})
	if len(got) != 1 || got[0] != 443 {
		t.Fatalf("expected [443], got %v", got)
	}
}

func TestChain_MultipleFilters(t *testing.T) {
	f1 := filter.New(filter.Rule{ExcludePorts: []int{80}})
	f2 := filter.New(filter.Rule{ExcludeRanges: []filter.Range{{Low: 400, High: 500}}})
	c := filter.NewChain(f1, f2)
	got := c.Apply([]int{80, 443, 8080})
	if len(got) != 1 || got[0] != 8080 {
		t.Fatalf("expected [8080], got %v", got)
	}
}

func TestChain_Len(t *testing.T) {
	f1 := filter.New(filter.Rule{})
	f2 := filter.New(filter.Rule{})
	c := filter.NewChain(f1, f2)
	if c.Len() != 2 {
		t.Fatalf("expected len 2, got %d", c.Len())
	}
}
