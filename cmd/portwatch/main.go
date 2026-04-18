package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/state"
)

func main() {
	var (
		portRange = flag.String("ports", "1-1024", "Port range to scan (e.g. 1-1024)")
		stateFile = flag.String("state", "/var/lib/portwatch/state.json", "Path to state file")
		alertEmail = flag.String("email", "", "Email address to send alerts (optional)")
		workers   = flag.Int("workers", 100, "Number of concurrent scan workers")
	)
	flag.Parse()

	s, err := scanner.New(*portRange, *workers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid port range: %v\n", err)
		os.Exit(1)
	}

	current, err := s.Scan()
	if err != nil {
		log.Fatalf("scan failed: %v", err)
	}

	previous, err := state.Load(*stateFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to load state: %v", err)
	}

	diff := state.Diff(previous, current)

	if err := state.Save(*stateFile, current); err != nil {
		log.Fatalf("failed to save state: %v", err)
	}

	a, err := alert.New(*alertEmail)
	if err != nil {
		log.Fatalf("failed to create alerter: %v", err)
	}

	if err := a.Notify(diff); err != nil {
		log.Fatalf("failed to send alert: %v", err)
	}
}
