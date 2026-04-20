// Package portwatch provides the core pipeline, runner, hooks, and
// lifecycle management for portwatch.
//
// # Lifecycle
//
// Lifecycle tracks named stages of a watcher run and broadcasts them to
// registered handlers. Supported events:
//
//   - starting  – emitted before any work begins.
//   - ready     – emitted once initialisation is complete.
//   - scan_begin – emitted at the start of each port scan.
//   - scan_end   – emitted when a port scan finishes.
//   - stopping  – emitted when a shutdown signal is received.
//   - stopped   – emitted after all cleanup is complete.
//
// Built-in middleware:
//
//   - ScanTimingMiddleware – logs the duration of each scan.
//   - UptimeMiddleware     – logs total uptime when the watcher stops.
package portwatch
