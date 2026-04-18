// Package watcher ties together scanning, state diffing, and reporting
// into a single reusable watch cycle.
package watcher

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Watcher performs a single scan-diff-report cycle.
type Watcher struct {
	scanner  *scanner.Scanner
	statePath string
	reporter *reporter.Reporter
}

// Config holds the dependencies needed to build a Watcher.
type Config struct {
	Scanner   *scanner.Scanner
	StatePath string
	Writer    io.Writer
	Format    string
}

// New creates a Watcher from the provided Config.
func New(cfg Config) (*Watcher, error) {
	if cfg.Scanner == nil {
		return nil, fmt.Errorf("watcher: scanner must not be nil")
	}
	if cfg.StatePath == "" {
		return nil, fmt.Errorf("watcher: state path must not be empty")
	}
	w := cfg.Writer
	if w == nil {
		w = os.Stdout
	}
	r := reporter.New(cfg.Format, w)
	return &Watcher{
		scanner:   cfg.Scanner,
		statePath: cfg.StatePath,
		reporter:  r,
	}, nil
}

// Run executes one watch cycle: scan ports, load previous state, diff, report, save.
func (w *Watcher) Run(ctx context.Context) error {
	ports, err := w.scanner.Scan(ctx)
	if err != nil {
		return fmt.Errorf("watcher: scan failed: %w", err)
	}

	prev, err := state.Load(w.statePath)
	if err != nil {
		return fmt.Errorf("watcher: load state failed: %w", err)
	}

	diff := state.Diff(prev, ports)

	if err := w.reporter.Report(diff); err != nil {
		return fmt.Errorf("watcher: report failed: %w", err)
	}

	if err := state.Save(w.statePath, ports); err != nil {
		return fmt.Errorf("watcher: save state failed: %w", err)
	}

	return nil
}
