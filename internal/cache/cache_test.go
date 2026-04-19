package cache_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/cache"
)

func TestGet_MissingKey(t *testing.T) {
	c := cache.New(time.Minute)
	_, ok := c.Get("localhost")
	if ok {
		t.Fatal("expected miss for unknown key")
	}
}

func TestSetAndGet(t *testing.T) {
	c := cache.New(time.Minute)
	ports := []int{80, 443, 8080}
	c.Set("localhost", ports)
	got, ok := c.Get("localhost")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(got) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(got))
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("localhost", []int{22})
	time.Sleep(20 * time.Millisecond)
	_, ok := c.Get("localhost")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestInvalidate(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("localhost", []int{80})
	c.Invalidate("localhost")
	_, ok := c.Get("localhost")
	if ok {
		t.Fatal("expected miss after invalidation")
	}
}

func TestFlush(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("host1", []int{80})
	c.Set("host2", []int{443})
	c.Flush()
	_, ok1 := c.Get("host1")
	_, ok2 := c.Get("host2")
	if ok1 || ok2 {
		t.Fatal("expected all entries cleared after flush")
	}
}
