package portname

import (
	"testing"
)

func TestLookup_Builtin(t *testing.T) {
	r := New(nil)
	if got := r.Lookup(22); got != "ssh" {
		t.Fatalf("expected ssh, got %s", got)
	}
}

func TestLookup_Unknown(t *testing.T) {
	r := New(nil)
	if got := r.Lookup(9999); got != "port/9999" {
		t.Fatalf("expected port/9999, got %s", got)
	}
}

func TestLookup_ExtraOverridesBuiltin(t *testing.T) {
	r := New(map[int]string{22: "custom-ssh"})
	if got := r.Lookup(22); got != "custom-ssh" {
		t.Fatalf("expected custom-ssh, got %s", got)
	}
}

func TestLookup_ExtraNew(t *testing.T) {
	r := New(map[int]string{12345: "myapp"})
	if got := r.Lookup(12345); got != "myapp" {
		t.Fatalf("expected myapp, got %s", got)
	}
}

func TestAnnotate_Empty(t *testing.T) {
	r := New(nil)
	if got := r.Annotate(nil); len(got) != 0 {
		t.Fatalf("expected empty slice")
	}
}

func TestAnnotate_FormatsCorrectly(t *testing.T) {
	r := New(nil)
	got := r.Annotate([]int{80, 443})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0] != "80(http)" {
		t.Errorf("unexpected: %s", got[0])
	}
	if got[1] != "443(https)" {
		t.Errorf("unexpected: %s", got[1])
	}
}

func TestLookup_NilResolver(t *testing.T) {
	var r *Resolver
	if got := r.Lookup(80); got != "http" {
		t.Fatalf("expected http, got %s", got)
	}
}
