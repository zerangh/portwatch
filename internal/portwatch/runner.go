package portwatch

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// RunnerConfig controls how the Runner repeats Pipeline cycles.
type RunnerConfig struct {
	Interval time.Duration
	MaxRuns  int // 0 = unlimited
}

// Runner repeatedly executes a Pipeline on a fixed interval.
type Runner struct {
	pipeline *Pipeline
	cfg      RunnerConfig
	log      *log.Logger
}

// NewRunner constructs a Runner wrapping the given Pipeline.
func NewRunner(p *Pipeline, rc RunnerConfig, out io.Writer) (*Runner, error) {
	if p == nil {
		return nil, fmt.Errorf("runner: pipeline must not be nil")
	}
	if rc.Interval <= 0 {
		return nil, fmt.Errorf("runner: interval must be positive")
	}
	if out == nil {
		out = os.Stderr
	}
	return &Runner{
		pipeline: p,
		cfg:      rc,
		log:      log.New(out, "[runner] ", log.LstdFlags),
	}, nil
}

// Start blocks, running the pipeline every Interval until ctx is cancelled
// or MaxRuns is reached.
func (r *Runner) Start(ctx context.Context) error {
	runs := 0
	ticker := time.NewTicker(r.cfg.Interval)
	defer ticker.Stop()

	// Run immediately on start.
	if err := r.runOnce(ctx, &runs); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := r.runOnce(ctx, &runs); err != nil {
				return err
			}
		}
	}
}

func (r *Runner) runOnce(ctx context.Context, runs *int) error {
	*runs++
	r.log.Printf("starting run #%d", *runs)
	if err := r.pipeline.Run(ctx); err != nil {
		r.log.Printf("run #%d error: %v", *runs, err)
	}
	if r.cfg.MaxRuns > 0 && *runs >= r.cfg.MaxRuns {
		return fmt.Errorf("runner: reached max runs (%d)", r.cfg.MaxRuns)
	}
	return nil
}
