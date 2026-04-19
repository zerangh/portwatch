// Package eventbus implements a lightweight in-process publish/subscribe bus
// used by portwatch to broadcast port-change events to decoupled consumers.
//
// Usage:
//
//	b := eventbus.New()
//	b.Subscribe("ports", func(e eventbus.Event) {
//		fmt.Println("change detected", e.Topic)
//	})
//	b.Publish(eventbus.Event{Topic: "ports", Snapshot: snap})
//
// Handlers are called synchronously in the goroutine that calls Publish.
// Wrap handlers with LoggingMiddleware or RecoveryMiddleware as needed.
package eventbus
