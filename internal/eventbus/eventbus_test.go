package eventbus_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/user/portwatch/internal/eventbus"
)

func TestSubscribeAndPublish(t *testing.T) {
	b := eventbus.New()
	var called int32
	b.Subscribe("ports", func(e eventbus.Event) {
		atomic.AddInt32(&called, 1)
	})
	b.Publish(eventbus.Event{Topic: "ports"})
	if atomic.LoadInt32(&called) != 1 {
		t.Fatalf("expected handler called once, got %d", called)
	}
}

func TestPublish_NoSubscribers(t *testing.T) {
	b := eventbus.New()
	// should not panic
	b.Publish(eventbus.Event{Topic: "ports"})
}

func TestUnsubscribe_RemovesHandlers(t *testing.T) {
	b := eventbus.New()
	b.Subscribe("ports", func(e eventbus.Event) {})
	b.Unsubscribe("ports")
	if b.Len("ports") != 0 {
		t.Fatal("expected no handlers after unsubscribe")
	}
}

func TestPublish_MultipleHandlers(t *testing.T) {
	b := eventbus.New()
	var count int32
	for i := 0; i < 3; i++ {
		b.Subscribe("ports", func(e eventbus.Event) {
			atomic.AddInt32(&count, 1)
		})
	}
	b.Publish(eventbus.Event{Topic: "ports"})
	if atomic.LoadInt32(&count) != 3 {
		t.Fatalf("expected 3 calls, got %d", count)
	}
}

func TestPublish_Concurrent(t *testing.T) {
	b := eventbus.New()
	var count int32
	b.Subscribe("ports", func(e eventbus.Event) {
		atomic.AddInt32(&count, 1)
	})
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Publish(eventbus.Event{Topic: "ports"})
		}()
	}
	wg.Wait()
	if atomic.LoadInt32(&count) != 20 {
		t.Fatalf("expected 20 calls, got %d", count)
	}
}
