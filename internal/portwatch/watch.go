// Package portwatch wires together the core scan-diff-alert pipeline.
package portwatch

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
	"github.com/user/portwatch/internal/state"
)

// Pipeline holds the components needed to run one watch cycle.
type Pipeline struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	alerter *alert.Alerter
	metrics *metrics.Metrics
	statePath string
	out     io.Writer
}

// NewPipeline constructs a Pipeline from the supplied config.
func NewPipeline(cfg *config.Config, statePath string, out io.Writer) (*Pipeline, error) {
	if cfg == nil {
		return nil, fmt.Errorf("portwatch: config must not be nil")
	}
	if statePath == "" {
		return nil, fmt.Errorf("portwatch: statePath must not be empty")
	}
	if out == nil {
		out = os.Stdout
	}
	sc, err := scanner.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("portwatch: scanner: %w", err)
	}
	al, err := alert.New(cfg, out)
	if err != nil {
		return nil, fmt.Errorf("portwatch: alert: %w", err)
	}
	return &Pipeline{
		cfg:       cfg,
		scanner:   sc,
		alerter:   al,
		metrics:   metrics.New(),
		statePath: statePath,
		out:       out,
	}, nil
}

// Run executes a single scan cycle: scan → diff → alert → save.
func (p *Pipeline) Run(ctx context.Context) error {
	ports, err := p.scanner.Scan(ctx)
	if err != nil {
		p.metrics.Record(nil, err)
		return fmt.Errorf("portwatch: scan: %w", err)
	}

	snap := snapshot.New(ports)

	prev, _ := state.Load(p.statePath)
	diff := state.Diff(prev, snap.Ports)

	p.metrics.Record(diff, nil)

	if err := p.alerter.Notify(diff); err != nil {
		return fmt.Errorf("portwatch: alert: %w", err)
	}

	if err := state.Save(p.statePath, snap.Ports); err != nil {
		return fmt.Errorf("portwatch: save state: %w", err)
	}
	return nil
}

// Metrics returns the accumulated scan metrics.
func (p *Pipeline) Metrics() *metrics.Metrics { return p.metrics }
