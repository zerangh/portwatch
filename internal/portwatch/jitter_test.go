package portwatch

import (
	"testing"
	"time"
)

func TestNewJitter_DefaultFactor(t *testing.T) {
	j := NewJitter(time.Second, 0) // 0 should be clamped to 0.1
	if j.factor != 0.1 {
		t.Fatalf("expected factor 0.1, got %v", j.factor)
	}
}

func TestNewJitter_ClampsFactor(t *testing.T) {
	j := NewJitter(time.Second, 5.0)
	if j.factor != 1.0 {
		t.Fatalf("expected factor clamped to 1.0, got %v", j.factor)
	}
}

func TestNext_WithinBounds(t *testing.T) {
	base := time.Second
	factor := 0.2
	j := NewJitter(base, factor)

	max := base + time.Duration(float64(base)*factor)
	min := base - time.Duration(float64(base)*factor)
	if min < time.Millisecond {
		min = time.Millisecond
	}

	for i := 0; i < 200; i++ {
		d := j.Next()
		if d < min || d > max {
			t.Fatalf("Next() = %v out of [%v, %v]", d, min, max)
		}
	}
}

func TestNext_NeverBelowMillisecond(t *testing.T) {
	j := NewJitter(time.Millisecond, 1.0) // factor=1 means offset can be -1ms
	for i := 0; i < 100; i++ {
		if d := j.Next(); d < time.Millisecond {
			t.Fatalf("Next() = %v, want >= 1ms", d)
		}
	}
}

func TestReset_ChangesBase(t *testing.T) {
	j := NewJitter(time.Second, 0.1)
	j.Reset(5 * time.Second)
	if j.Base() != 5*time.Second {
		t.Fatalf("Base() = %v, want 5s", j.Base())
	}
}

func TestBase_ReturnsBase(t *testing.T) {
	j := NewJitter(3*time.Second, 0.1)
	if j.Base() != 3*time.Second {
		t.Fatalf("Base() = %v, want 3s", j.Base())
	}
}
