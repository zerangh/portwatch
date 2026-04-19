// Package debounce provides a Debouncer that coalesces rapid successive
// function calls into a single invocation after a configurable quiet period.
//
// This is useful in portwatch to avoid triggering repeated alerts or state
// saves when multiple port changes are detected in quick succession during
// a scan interval.
//
// Example usage:
//
//	d := debounce.New(500 * time.Millisecond)
//	// Each call resets the timer; fn fires once after 500ms of quiet.
//	d.Call(fn)
package debounce
