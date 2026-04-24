package portwatch

import (
	"testing"
	"time"
)

func TestCooldown_FirstCallAllowed(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	if !c.Allow("key1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestCooldown_SecondCallBlocked(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	c.Allow("key1")
	if c.Allow("key1") {
		t.Fatal("expected second call within window to be blocked")
	}
}

func TestCooldown_AllowsAfterWindowExpires(t *testing.T) {
	now := time.Now()
	c := NewCooldown(5 * time.Second)
	c.now = func() time.Time { return now }
	c.Allow("key1")

	// advance time past the window
	c.now = func() time.Time { return now.Add(6 * time.Second) }
	if !c.Allow("key1") {
		t.Fatal("expected call after window expiry to be allowed")
	}
}

func TestCooldown_ZeroWindowAlwaysAllows(t *testing.T) {
	c := NewCooldown(0)
	for i := 0; i < 5; i++ {
		if !c.Allow("key") {
			t.Fatalf("iteration %d: expected allow with zero window", i)
		}
	}
}

func TestCooldown_IndependentKeys(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	c.Allow("a")
	if !c.Allow("b") {
		t.Fatal("different key should be allowed")
	}
}

func TestCooldown_Reset_AllowsImmediately(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	c.Allow("key1")
	c.Reset("key1")
	if !c.Allow("key1") {
		t.Fatal("expected allow after reset")
	}
}

func TestCooldown_Flush_RemovesExpired(t *testing.T) {
	now := time.Now()
	c := NewCooldown(5 * time.Second)
	c.now = func() time.Time { return now }
	c.Allow("a")
	c.Allow("b")

	c.now = func() time.Time { return now.Add(6 * time.Second) }
	c.Flush()

	if c.Len() != 0 {
		t.Fatalf("expected 0 entries after flush, got %d", c.Len())
	}
}

func TestCooldown_Flush_KeepsActive(t *testing.T) {
	now := time.Now()
	c := NewCooldown(10 * time.Second)
	c.now = func() time.Time { return now }
	c.Allow("a")

	c.now = func() time.Time { return now.Add(3 * time.Second) }
	c.Flush()

	if c.Len() != 1 {
		t.Fatalf("expected 1 active entry after flush, got %d", c.Len())
	}
}
