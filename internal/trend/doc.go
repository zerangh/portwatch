// Package trend provides lightweight analysis of port-change trends
// derived from scan history.
//
// Use Analyze to compute a Result from a slice of history.Entry values,
// optionally restricting the window with a since duration.
//
// Use Print to render a human-readable summary to any io.Writer.
package trend
