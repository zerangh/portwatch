// Package portwatch provides the core pipeline, runner, hooks, and lifecycle
// management for portwatch.
//
// # Cooldown
//
// The Cooldown type suppresses repeated alert events for the same port-change
// key within a configurable time window. This prevents alert storms when a
// port flaps rapidly or when a scan runs more frequently than alerts can be
// meaningfully acted upon.
//
// Usage:
//
//	cd := portwatch.NewCooldown(30 * time.Second)
//
//	if cd.Allow("added:8080") {
//		// send alert
//	}
//
// Call Reset to immediately unblock a specific key, or Flush to evict all
// entries whose window has already expired (useful in long-running daemons to
// keep memory bounded).
package portwatch
