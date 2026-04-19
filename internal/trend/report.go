package trend

import (
	"fmt"
	"io"
	"os"
)

// Print writes a human-readable trend summary to w.
func Print(r Result, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	if r.Samples == 0 {
		fmt.Fprintln(w, "trend: no data")
		return
	}
	fmt.Fprintf(w, "trend: %s (net %+d over %d scans, +%.2f/scan -%.2f/scan)\n",
		r.Direction, r.NetChange, r.Samples, r.AddedRate, r.RemovedRate)
}
