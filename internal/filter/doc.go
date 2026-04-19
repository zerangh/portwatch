// Package filter provides composable port filtering for portwatch.
//
// A Filter is constructed with a Rule that specifies ports or port ranges
// to exclude from scan results. This allows operators to suppress
// well-known or expected ports from change detection.
//
// Example:
//
//	f := filter.New(filter.Rule{
//		ExcludePorts:  []int{22, 80, 443},
//		ExcludeRanges: []filter.Range{{Low: 32768, High: 60999}},
//	})
//	filtered := f.Apply(scannedPorts)
package filter
