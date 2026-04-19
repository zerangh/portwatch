// Package healthcheck tracks the operational health of portwatch scan cycles.
//
// A Checker accumulates scan and error events, exposing a Status summary
// and a human-readable Print method suitable for CLI output or log lines.
//
// Example:
//
//	ch := healthcheck.New(os.Stdout)
//	ch.RecordScan(time.Now())
//	ch.Print()
package healthcheck
