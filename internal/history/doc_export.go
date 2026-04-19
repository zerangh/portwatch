// Package history provides scan history storage, querying, and export.
//
// Export
//
// ExportJSON writes history entries as a JSON array to the given writer.
// ExportCSV writes history entries in CSV format with a header row.
// Export dispatches to the correct format based on a format string ("json" or "csv").
// Returns an error if the format string is not recognized.
//
// Supported format values:
//   - "json": JSON array output
//   - "csv":  CSV output with a header row
//
// Example:
//
//	f, _ := os.Create("history.csv")
//	defer f.Close()
//	if err := history.Export(entries, "csv", f); err != nil {
//		log.Fatal(err)
//	}
package history
