package portgroup_test

import (
	"testing"

	"github.com/user/portwatch/internal/portgroup"
)

func TestDefineAndLookup(t *testing.T) {
	r := portgroup.New()
	r.Define("web", []int{443, 80, 8080})
	ports, ok := r.Lookup("web")
	if !ok {
		t.Fatal("expected group 'web' to exist")
	}
	if len(ports) != 3 || ports[0] != 80 {
		t.Fatalf("expected sorted ports [80 443 8080], got %v", ports)
	}
}

func TestLookup_Missing(t *testing.T) {
	r := portgroup.New()
	_, ok := r.Lookup("nope")
	if ok {
		t.Fatal("expected missing group to return false")
	}
}

func TestClassify_SingleGroup(t *testing.T) {
	r := portgroup.New()
	r.Define("db", []int{5432, 3306})
	names := r.Classify(5432)
	if len(names) != 1 || names[0] != "db" {
		t.Fatalf("unexpected classify result: %v", names)
	}
}

func TestClassify_MultipleGroups(t *testing.T) {
	r := portgroup.New()
	r.Define("web", []int{80, 443})
	r.Define("all", []int{80, 443, 5432})
	names := r.Classify(80)
	if len(names) != 2 {
		t.Fatalf("expected 2 groups, got %v", names)
	}
}

func TestClassify_NoMatch(t *testing.T) {
	r := portgroup.New()
	r.Define("web", []int{80})
	names := r.Classify(9999)
	if len(names) != 0 {
		t.Fatalf("expected no groups, got %v", names)
	}
}

func TestAll_SortedByName(t *testing.T) {
	r := portgroup.New()
	r.Define("z", []int{1})
	r.Define("a", []int{2})
	groups := r.All()
	if len(groups) != 2 || groups[0].Name != "a" {
		t.Fatalf("expected sorted groups, got %v", groups)
	}
}
