// Package reporter provides formatted output of port scan diff results.
//
// It supports multiple output formats (text, JSON) and can write to any
// io.Writer, defaulting to os.Stdout.
//
// Usage:
//
//	r := reporter.New(reporter.FormatText, nil)
//	r.Report(diff)
package reporter
