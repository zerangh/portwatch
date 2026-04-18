// Package notifier defines the Notifier interface and built-in backends
// for delivering port-change alerts in portwatch.
//
// Backends:
//
//   - StdoutNotifier — writes human-readable messages to any io.Writer
//     (defaults to os.Stdout). Suitable for terminal use and log pipelines.
//
// Additional backends (e.g. webhook, email) can be added by implementing
// the Notifier interface:
//
//	type Notifier interface {
//	    Notify(added, removed []state.Port) error
//	}
package notifier
