package resolve

import (
	"testing"
	"time"
)

func TestNew_DefaultTimeout(t *testing.T) {
	r := New(0)
	if r.timeout != 5*time.Second {
		t.Fatalf("expected 5s default timeout, got %v", r.timeout)
	}
}

func TestResolve_EmptyHost(t *testing.T) {
	r := New(5 * time.Second)
	_, err := r.Resolve("")
	if err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestResolve_IPAddress(t *testing.T) {
	r := New(5 * time.Second)
	res, err := r.Resolve("127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Addresses) != 1 || res.Addresses[0] != "127.0.0.1" {
		t.Fatalf("expected [127.0.0.1], got %v", res.Addresses)
	}
	if res.Host != "127.0.0.1" {
		t.Fatalf("expected host 127.0.0.1, got %q", res.Host)
	}
}

func TestResolve_Localhost(t *testing.T) {
	r := New(5 * time.Second)
	res, err := r.Resolve("localhost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Addresses) == 0 {
		t.Fatal("expected at least one address for localhost")
	}
	if res.ResolvedAt.IsZero() {
		t.Fatal("expected non-zero ResolvedAt")
	}
}

func TestPrimary_NoAddresses(t *testing.T) {
	res := &Result{Host: "example", Addresses: []string{}}
	_, err := res.Primary()
	if err == nil {
		t.Fatal("expected error when no addresses")
	}
}

func TestPrimary_ReturnsFirst(t *testing.T) {
	res := &Result{Host: "example", Addresses: []string{"1.2.3.4", "5.6.7.8"}}
	ip, err := res.Primary()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip != "1.2.3.4" {
		t.Fatalf("expected 1.2.3.4, got %q", ip)
	}
}
