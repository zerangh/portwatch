package portwatch

import (
	"math/rand"
	"time"
)

// Jitter adds randomised delay to a base interval to spread out scan
// bursts when multiple portwatch instances run on the same schedule.
type Jitter struct {
	base    time.Duration
	factor  float64 // fraction of base to use as max jitter, e.g. 0.2 = ±20%
	rng     *rand.Rand
}

// NewJitter returns a Jitter that randomises intervals by up to factor*base.
// factor must be in the range (0, 1]; values outside that range are clamped.
func NewJitter(base time.Duration, factor float64) *Jitter {
	if factor <= 0 {
		factor = 0.1
	}
	if factor > 1 {
		factor = 1
	}
	return &Jitter{
		base:   base,
		factor: factor,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Next returns base ± a random offset bounded by factor*base.
// The result is always >= 1 millisecond.
func (j *Jitter) Next() time.Duration {
	max := float64(j.base) * j.factor
	offset := time.Duration((j.rng.Float64()*2-1)*max) // range: [-max, +max)
	d := j.base + offset
	if d < time.Millisecond {
		d = time.Millisecond
	}
	return d
}

// Reset replaces the base interval, leaving the factor unchanged.
func (j *Jitter) Reset(base time.Duration) {
	j.base = base
}

// Base returns the current base interval.
func (j *Jitter) Base() time.Duration {
	return j.base
}
