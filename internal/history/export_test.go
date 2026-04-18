package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var testEntries = []Entry{
	{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Added:     []int{80, 443},
		Removed:   []int{},
	},
	{
		Timestamp: time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC),
		Added:     []int{},
		Removed:   []int{8080},
	},
}

func TestExportJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportJSON(testEntries, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2024-01-15T10:00:00Z") {
		t.Error("expected timestamp in output")
	}
	if !strings.Contains(out, "80") || !strings.Contains(out, "443") {
		t.Error("expected added ports in output")
	}
	if !strings.Contains(out, "8080") {
		t.Error("expected removed port in output")
	}
}

func TestExportCSV(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportCSV(testEntries, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
	if lines[0] != "timestamp,added,removed" {
		t.Errorf("unexpected header: %s", lines[0])
	}
	if !strings.Contains(lines[1], "80;443") {
		t.Errorf("expected ports in csv row: %s", lines[1])
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(testEntries, ExportFormat("xml"), &buf)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExport_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportJSON([]Entry{}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[]") {
		t.Error("expected empty JSON array")
	}
}
