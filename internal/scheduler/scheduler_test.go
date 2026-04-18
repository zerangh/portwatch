package scheduler_test

import (
	"context"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scheduler"
)

func TestScheduler_RunsJobImmediately(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.New(os.Discard, "", 0)
	s := scheduler.New(10*time.Second, job, logger)

	go s.Run(ctx)

	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&count) < 1 {
		t.Error("expected job to run immediately, but it did not")
	}
}

func TestScheduler_RunsJobOnInterval(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.New(os.Discard, "", 0)
	s := scheduler.New(50*time.Millisecond, job, logger)

	go s.Run(ctx)

	time.Sleep(180 * time.Millisecond)
	cancel()

	final := atomic.LoadInt32(&count)
	if final < 3 {
		t.Errorf("expected at least 3 job executions, got %d", final)
	}
}

func TestScheduler_StopsOnContextCancel(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	logger := log.New(os.Discard, "", 0)
	s := scheduler.New(20*time.Millisecond, job, logger)

	go s.Run(ctx)
	time.Sleep(30 * time.Millisecond)
	cancel()
	time.Sleep(60 * time.Millisecond)

	snapshot := atomic.LoadInt32(&count)
	time.Sleep(60 * time.Millisecond)
	if atomic.LoadInt32(&count) != snapshot {
		t.Error("scheduler continued running after context was cancelled")
	}
}
