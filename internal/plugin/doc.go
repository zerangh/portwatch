// Package plugin implements a lightweight hook-based extension system for
// portwatch. Callers register named Handler functions with a Registry; when a
// port-change event occurs the Registry dispatches the event to every
// registered handler.
//
// Built-in helpers:
//
//   - LogHandler – writes a human-readable line to any io.Writer.
//   - ThresholdHandler – gates another handler behind a minimum change count.
//
// Example:
//
//	reg := plugin.New()
//	reg.Register("log", plugin.LogHandler(os.Stdout))
//	reg.Dispatch(plugin.Event{Host: "localhost", Added: []int{8080}})
package plugin
