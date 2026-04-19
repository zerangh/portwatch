package plugin

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// LogHandler returns a Handler that writes a human-readable summary of each
// event to w. If w is nil, os.Stdout is used.
func LogHandler(w io.Writer) Handler {
	if w == nil {
		w = os.Stdout
	}
	return func(e Event) error {
		if len(e.Added) == 0 && len(e.Removed) == 0 {
			return nil
		}
		parts := []string{}
		if len(e.Added) > 0 {
			parts = append(parts, fmt.Sprintf("added=%v", e.Added))
		}
		if len(e.Removed) > 0 {
			parts = append(parts, fmt.Sprintf("removed=%v", e.Removed))
		}
		_, err := fmt.Fprintf(w, "[plugin:log] host=%s %s\n", e.Host, strings.Join(parts, " "))
		return err
	}
}

// ThresholdHandler returns a Handler that invokes action only when the total
// number of changed ports (added + removed) meets or exceeds threshold.
func ThresholdHandler(threshold int, action Handler) Handler {
	return func(e Event) error {
		if len(e.Added)+len(e.Removed) >= threshold {
			return action(e)
		}
		return nil
	}
}
