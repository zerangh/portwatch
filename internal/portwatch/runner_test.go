package portwatch_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portwatch"
)

func TestNewRunner_NilPipeline(t *testing.T) {
	_, err := portwatch.NewRunner(nil, portwatch.RunnerConfig{Interval: time.Second}, nil)
	if err == nil {
		t.Fatal("expected error for nil pipeline")
	}
}

func TestNewRunner_ZeroInterval(t *testing.T) {
	p, _ := portwatch.NewPipeline(defaultCfg(), tmpState(t), &bytes.Buffer{})
	_, err := portwatch.NewRunner(p, portwatch.RunnerConfig{Interval: 0}, nil)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestRunner_StopsAfterMaxRuns(t *testing.T) {
	sp := tmpState(t)
	p, err := portwatch.NewPipeline(defaultCfg(), sp, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("NewPipeline: %v", err)
	}
	var logBuf bytes.Buffer
	r, err := portwatch.NewRunner(p, portwatch.RunnerConfig{
		Interval: 10 * time.Millisecond,
		MaxRuns:  2,
	}, &logBuf)
	if err != nil {
		t.Fatalf("NewRunner: %v", err)
	}
	err = r.Start(context.Background())
	if err == nil || !strings.Contains(err.Error(), "max runs") {
		t.Fatalf("expected max-runs error, got: %v", err)
	}
}

func TestRunner_StopsOnContextCancel(t *testing.T) {
	sp := tmpState(t)
	p, err := portwatch.NewPipeline(defaultCfg(), sp, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("NewPipeline: %v", err)
	}
	r, err := portwatch.NewRunner(p, portwatch.RunnerConfig{
		Interval: 20 * time.Millisecond,
	}, nil)
	if err != nil {
		t.Fatalf("NewRunner: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Millisecond)
	defer cancel()
	err = r.Start(ctx)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}
