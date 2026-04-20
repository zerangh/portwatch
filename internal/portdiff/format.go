package portdiff

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Format controls how a Result is rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Print writes a formatted diff to w. If w is nil, os.Stdout is used.
// format must be FormatText or FormatJSON.
func Print(w io.Writer, r Result, format Format) error {
	if w == nil {
		w = os.Stdout
	}
	switch format {
	case FormatJSON:
		return printJSON(w, r)
	default:
		return printText(w, r)
	}
}

func printText(w io.Writer, r Result) error {
	if r.IsEmpty() {
		_, err := fmt.Fprintln(w, "no port changes detected")
		return err
	}
	var sb strings.Builder
	if len(r.Added) > 0 {
		sb.WriteString(fmt.Sprintf("added:   %s\n", joinInts(r.Added)))
	}
	if len(r.Removed) > 0 {
		sb.WriteString(fmt.Sprintf("removed: %s\n", joinInts(r.Removed)))
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}

func printJSON(w io.Writer, r Result) error {
	added := intSlice(r.Added)
	removed := intSlice(r.Removed)
	_, err := fmt.Fprintf(w, `{"added":%s,"removed":%s}\n`, added, removed)
	return err
}

func joinInts(nums []int) string {
	parts := make([]string, len(nums))
	for i, n := range nums {
		parts[i] = fmt.Sprintf("%d", n)
	}
	return strings.Join(parts, ", ")
}

func intSlice(nums []int) string {
	if len(nums) == 0 {
		return "[]"
	}
	return "[" + joinInts(nums) + "]"
}
