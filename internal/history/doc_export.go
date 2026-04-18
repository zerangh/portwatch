// Package history provides scan history storage, querying, and export.
//
// Export
//
// ExportJSON writes history entries as a JSON array to the given writer.
// ExportCSV writes history entries in CSV format with a header row.
// Export dispatches to the correct format based on a format string ("json" or "csv").
//
// Example:
//
//	f, _ := os.Create("history.csv")
//	defer f.Close()
//	history.Export(entries, "csv", f)
package history
