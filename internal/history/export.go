package history

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// ExportFormat defines the output format for history export.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
)

// ExportJSON writes all history entries as a JSON array to w.
func ExportJSON(entries []Entry, w io.Writer) error {
	type row struct {
		Timestamp string `json:"timestamp"`
		Added     []int  `json:"added"`
		Removed   []int  `json:"removed"`
	}
	rows := make([]row, len(entries))
	for i, e := range entries {
		rows[i] = row{
			Timestamp: e.Timestamp.Format(time.RFC3339),
			Added:     e.Added,
			Removed:   e.Removed,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}

// ExportCSV writes all history entries as CSV to w.
func ExportCSV(entries []Entry, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"timestamp", "added", "removed"}); err != nil {
		return err
	}
	for _, e := range entries {
		cw.Write([]string{
			e.Timestamp.Format(time.RFC3339),
			joinInts(e.Added),
			joinInts(e.Removed),
		})
	}
	cw.Flush()
	return cw.Error()
}

// Export dispatches to the appropriate exporter based on format.
func Export(entries []Entry, format ExportFormat, w io.Writer) error {
	switch format {
	case FormatJSON:
		return ExportJSON(entries, w)
	case FormatCSV:
		return ExportCSV(entries, w)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func joinInts(vals []int) string {
	if len(vals) == 0 {
		return ""
	}
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ";")
}
