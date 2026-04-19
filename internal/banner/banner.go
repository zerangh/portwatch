// Package banner prints a startup summary to the terminal when portwatch begins.
package banner

import (
	"fmt"
	"io"
	"os"
	"time"
)

const version = "0.1.0"

// Options controls what is shown in the banner.
type Options struct {
	Host      string
	PortRange string
	Interval  time.Duration
	StatePath string
	Writer    io.Writer
}

// Print writes the startup banner to the configured writer (defaults to os.Stdout).
func Print(opts Options) {
	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}

	fmt.Fprintf(w, "╔══════════════════════════════════╗\n")
	fmt.Fprintf(w, "║   portwatch v%-20s║\n", version+"   ")
	fmt.Fprintf(w, "╚══════════════════════════════════╝\n")
	fmt.Fprintf(w, "  host       : %s\n", opts.Host)
	fmt.Fprintf(w, "  port range : %s\n", opts.PortRange)
	fmt.Fprintf(w, "  interval   : %s\n", opts.Interval)
	fmt.Fprintf(w, "  state file : %s\n", opts.StatePath)
	fmt.Fprintf(w, "  started at : %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintln(w)
}
