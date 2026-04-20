package portwatch_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/portwatch"
)

func tmpState(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func defaultCfg() *config.Config {
	cfg := config.DefaultConfig()
	// Use a tiny port range unlikely to be open in CI.
	cfg.PortRange = "65530-65534"
	cfg.Concurrency = 2
	return cfg
}

func TestNewPipeline_NilConfig(t *testing.T) {
	_, err := portwatch.NewPipeline(nil, tmpState(t), nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestNewPipeline_EmptyStatePath(t *testing.T) {
	_, err := portwatch.NewPipeline(defaultCfg(), "", nil)
	if err == nil {
		t.Fatal("expected error for empty state path")
	}
}

func TestNewPipeline_NilWriterUsesStdout(t *testing.T) {
	p, err := portwatch.NewPipeline(defaultCfg(), tmpState(t), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestRun_CreatesStateFile(t *testing.T) {
	sp := tmpState(t)
	var buf bytes.Buffer
	p, err := portwatch.NewPipeline(defaultCfg(), sp, &buf)
	if err != nil {
		t.Fatalf("NewPipeline: %v", err)
	}
	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if _, err := os.Stat(sp); err != nil {
		t.Fatalf("state file not created: %v", err)
	}
}

func TestRun_MetricsRecorded(t *testing.T) {
	sp := tmpState(t)
	p, err := portwatch.NewPipeline(defaultCfg(), sp, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("NewPipeline: %v", err)
	}
	_ = p.Run(context.Background())
	m := p.Metrics()
	if m == nil {
		t.Fatal("expected non-nil metrics")
	}
}
