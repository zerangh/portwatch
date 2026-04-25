package portwatch

import (
	"io"
	"os"
	"sync"
	"time"
)

// Sampler records periodic scan durations and computes rolling statistics
// such as mean and p95 latency over a sliding window of samples.
type Sampler struct {
	mu      sync.Mutex
	samples []time.Duration
	maxSize int
	w       io.Writer
}

// NewSampler returns a Sampler that retains up to maxSize duration samples.
// If maxSize is zero it defaults to 100. A nil writer falls back to os.Stdout.
func NewSampler(maxSize int, w io.Writer) *Sampler {
	if maxSize <= 0 {
		maxSize = 100
	}
	if w == nil {
		w = os.Stdout
	}
	return &Sampler{maxSize: maxSize, w: w}
}

// Record appends d to the sample window, evicting the oldest entry when full.
func (s *Sampler) Record(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.samples) >= s.maxSize {
		s.samples = s.samples[1:]
	}
	s.samples = append(s.samples, d)
}

// Mean returns the arithmetic mean of all recorded samples.
// It returns 0 when no samples have been recorded.
func (s *Sampler) Mean() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.samples) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range s.samples {
		total += d
	}
	return total / time.Duration(len(s.samples))
}

// P95 returns the 95th-percentile duration from the recorded samples.
// It returns 0 when no samples have been recorded.
func (s *Sampler) P95() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := len(s.samples)
	if n == 0 {
		return 0
	}
	sorted := make([]time.Duration, n)
	copy(sorted, s.samples)
	sortDurations(sorted)
	idx := int(float64(n)*0.95) - 1
	if idx < 0 {
		idx = 0
	}
	return sorted[idx]
}

// Len returns the current number of recorded samples.
func (s *Sampler) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.samples)
}

func sortDurations(d []time.Duration) {
	for i := 1; i < len(d); i++ {
		for j := i; j > 0 && d[j] < d[j-1]; j-- {
			d[j], d[j-1] = d[j-1], d[j]
		}
	}
}
