// Package alert provides notification mechanisms for port change events.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Notifier sends alerts about port changes.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// Pass nil to use os.Stdout.
func New(out io.Writer) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	return &Notifier{out: out}
}

// Notify prints a human-readable summary of the diff to the writer.
// Returns true if any changes were detected.
func (n *Notifier) Notify(diff state.Diff) bool {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return false
	}

	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(n.out, "[%s] Port changes detected:\n", timestamp)

	for _, p := range diff.Opened {
		fmt.Fprintf(n.out, "  + OPENED  %s\n", p)
	}
	for _, p := range diff.Closed {
		fmt.Fprintf(n.out, "  - CLOSED  %s\n", p)
	}

	return true
}
