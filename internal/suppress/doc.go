// Package suppress manages a persistent suppression list of ports that should
// be silently ignored during change detection and alerting.
//
// Ports added to the suppression list are excluded from diff output and will
// never trigger notifications, even if their state changes between scans.
//
// The list is stored as a JSON file and is safe for concurrent use.
//
// Example:
//
//	l, err := suppress.New("/var/lib/portwatch/suppress.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	l.Add(22)   // never alert on SSH
//	visible := l.Filter(openPorts)
package suppress
