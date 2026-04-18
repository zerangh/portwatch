// Package reporter provides structured reporting of port scan results.
package reporter

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes port change reports to an output destination.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a new Reporter with the given format and output writer.
// If out is nil, os.Stdout is used.
func New(format Format, out io.Writer) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out, format: format}
}

// Report writes a formatted report of the given diff to the output.
func (r *Reporter) Report(diff state.Diff) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(diff)
	default:
		return r.writeText(diff)
	}
}

func (r *Reporter) writeText(diff state.Diff) error {
	ts := time.Now().Format(time.RFC3339)
	if len(diff.Added) == 0 && len(diff.Removed) == 0 {
		_, err := fmt.Fprintf(r.out, "[%s] No port changes detected.\n", ts)
		return err
	}
	for _, p := range diff.Added {
		if _, err := fmt.Fprintf(r.out, "[%s] OPENED  port %d\n", ts, p); err != nil {
			return err
		}
	}
	for _, p := range diff.Removed {
		if _, err := fmt.Fprintf(r.out, "[%s] CLOSED  port %d\n", ts, p); err != nil {
			return err
		}
	}
	return nil
}

func (r *Reporter) writeJSON(diff state.Diff) error {
	ts := time.Now().Format(time.RFC3339)
	_, err := fmt.Fprintf(r.out,
		`{"timestamp":%q,"added":%s,"removed":%s}`+"\n",
		ts, intSliceJSON(diff.Added), intSliceJSON(diff.Removed))
	return err
}

func intSliceJSON(ports []int) string {
	if len(ports) == 0 {
		return "[]"
	}
	out := "["
	for i, p := range ports {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%d", p)
	}
	return out + "]"
}
