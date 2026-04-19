// Package plugin provides a simple hook-based plugin system for portwatch.
// Plugins can register callbacks that are invoked when port changes are detected.
package plugin

import "sync"

// Event represents a port change event passed to plugins.
type Event struct {
	Added   []int
	Removed []int
	Host    string
}

// Handler is a function that handles a plugin event.
type Handler func(e Event) error

// Registry holds registered plugin handlers.
type Registry struct {
	mu       sync.RWMutex
	handlers map[string]Handler
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{handlers: make(map[string]Handler)}
}

// Register adds a named handler to the registry.
// Registering the same name twice overwrites the previous handler.
func (r *Registry) Register(name string, h Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[name] = h
}

// Unregister removes a handler by name. It is a no-op if the name is absent.
func (r *Registry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.handlers, name)
}

// Len returns the number of registered handlers.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.handlers)
}

// Dispatch calls every registered handler with the given event.
// Errors are collected and returned as a slice; execution continues even if a
// handler returns an error.
func (r *Registry) Dispatch(e Event) []error {
	r.mu.RLock()
	handlers := make(map[string]Handler, len(r.handlers))
	for k, v := range r.handlers {
		handlers[k] = v
	}
	r.mu.RUnlock()

	var errs []error
	for _, h := range handlers {
		if err := h(e); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
