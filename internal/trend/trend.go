// Package trend analyses port change frequency over time.
package trend

import (
	"time"

	"github.com/user/portwatch/internal/history"
)

// Direction indicates whether open ports are growing, shrinking, or stable.
type Direction string

const (
	Growing  Direction = "growing"
	Shrinking Direction = "shrinking"
	Stable    Direction = "stable"
)

// Result holds a trend analysis result.
type Result struct {
	Direction  Direction
	NetChange  int
	AddedRate  float64 // average ports added per scan
	RemovedRate float64 // average ports removed per scan
	Samples    int
}

// Analyze computes a trend from the given history entries.
// Only entries within the since window are considered.
func Analyze(entries []history.Entry, since time.Duration) Result {
	cutoff := time.Now().Add(-since)
	var totalAdded, totalRemoved, samples int
	for _, e := range entries {
		if e.Timestamp.Before(cutoff) {
			continue
		}
		totalAdded += len(e.Diff.Added)
		totalRemoved += len(e.Diff.Removed)
		samples++
	}
	if samples == 0 {
		return Result{Direction: Stable, Samples: 0}
	}
	addedRate := float64(totalAdded) / float64(samples)
	removedRate := float64(totalRemoved) / float64(samples)
	net := totalAdded - totalRemoved
	dir := Stable
	switch {
	case net > 0:
		dir = Growing
	case net < 0:
		dir = Shrinking
	}
	return Result{
		Direction:   dir,
		NetChange:   net,
		AddedRate:   addedRate,
		RemovedRate: removedRate,
		Samples:     samples,
	}
}
