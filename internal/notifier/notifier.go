// Package notifier provides pluggable notification backends for portwatch.
package notifier

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/portwatch/internal/state"
)

// Notifier sends port change notifications via a backend.
type Notifier interface {
	Notify(added, removed []state.Port) error
}

// StdoutNotifier writes notifications to an io.Writer.
type StdoutNotifier struct {
	w io.Writer
}

// NewStdout returns a StdoutNotifier writing to w.
// If w is nil, os.Stdout is used.
func NewStdout(w io.Writer) *StdoutNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &StdoutNotifier{w: w}
}

// Notify writes a human-readable summary of port changes.
func (n *StdoutNotifier) Notify(added, removed []state.Port) error {
	if len(added) == 0 && len(removed) == 0 {
		return nil
	}
	var sb strings.Builder
	if len(added) > 0 {
		sb.WriteString(fmt.Sprintf("[portwatch] NEW ports open: %v\n", portList(added)))
	}
	if len(removed) > 0 {
		sb.WriteString(fmt.Sprintf("[portwatch] Ports closed: %v\n", portList(removed)))
	}
	_, err := fmt.Fprint(n.w, sb.String())
	return err
}

func portList(ports []state.Port) string {
	parts := make([]string, len(ports))
	for i, p := range ports {
		parts[i] = p.String()
	}
	return strings.Join(parts, ", ")
}
