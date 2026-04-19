// Package eventbus provides a simple publish/subscribe event bus for
// broadcasting port change events to multiple subscribers.
package eventbus

import (
	"sync"

	"github.com/user/portwatch/internal/snapshot"
)

// Event represents a port change event published on the bus.
type Event struct {
	Topic    string
	Snapshot *snapshot.Snapshot
	Prev     *snapshot.Snapshot
}

// Handler is a function that handles an Event.
type Handler func(Event)

// Bus is a simple in-process publish/subscribe event bus.
type Bus struct {
	mu       sync.RWMutex
	subs     map[string][]Handler
}

// New returns an initialised Bus.
func New() *Bus {
	return &Bus{subs: make(map[string][]Handler)}
}

// Subscribe registers h to receive events on topic.
func (b *Bus) Subscribe(topic string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[topic] = append(b.subs[topic], h)
}

// Unsubscribe removes all handlers for topic.
func (b *Bus) Unsubscribe(topic string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.subs, topic)
}

// Publish sends e to every handler subscribed to e.Topic.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.subs[e.Topic]))
	copy(handlers, b.subs[e.Topic])
	b.mu.RUnlock()
	for _, h := range handlers {
		h(e)
	}
}

// Len returns the number of handlers registered for topic.
func (b *Bus) Len(topic string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subs[topic])
}
