// Package baseline manages a known-good snapshot of open ports.
//
// A baseline represents a user-approved set of ports that are expected to be
// open on the host. When a baseline is active, portwatch suppresses alerts for
// ports that match the baseline, reporting only newly discovered or unexpected
// ports.
//
// Usage:
//
//	b := baseline.New("/var/lib/portwatch/baseline.json")
//	if err := b.Load(); err != nil {
//	    log.Fatal(err)
//	}
//	// Filter out known ports before alerting.
//	unexpected := b.Filter(scannedPorts)
package baseline
