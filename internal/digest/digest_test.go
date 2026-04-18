package digest

import (
	"testing"
)

func TestFromPorts_Empty(t *testing.T) {
	d := FromPorts(nil)
	if d != Empty {
		t.Errorf("expected Empty digest, got %s", d)
	}

	d2 := FromPorts([]int{})
	if d2 != Empty {
		t.Errorf("expected Empty digest for empty slice, got %s", d2)
	}
}

func TestFromPorts_Stable(t *testing.T) {
	a := FromPorts([]int{80, 443, 8080})
	b := FromPorts([]int{8080, 80, 443})
	if a != b {
		t.Errorf("expected same digest regardless of order: %s vs %s", a, b)
	}
}

func TestFromPorts_Different(t *testing.T) {
	a := FromPorts([]int{80, 443})
	b := FromPorts([]int{80, 444})
	if Equal(a, b) {
		t.Error("expected different digests for different port sets")
	}
}

func TestFromPorts_SinglePort(t *testing.T) {
	d := FromPorts([]int{22})
	if d == Empty {
		t.Error("single port should not produce empty digest")
	}
	if len(string(d)) == 0 {
		t.Error("digest should not be empty string")
	}
}

func TestEqual(t *testing.T) {
	a := FromPorts([]int{80})
	b := FromPorts([]int{80})
	if !Equal(a, b) {
		t.Error("expected Equal to return true for same ports")
	}
}

func TestDigest_String(t *testing.T) {
	d := FromPorts([]int{80})
	s := d.String()
	if len(s) == 0 {
		t.Error("String() should return non-empty value")
	}
	if s[:7] != "digest:" {
		t.Errorf("expected prefix 'digest:', got %s", s)
	}
}
