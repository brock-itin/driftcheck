package drift

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

var exportTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func exportFinding(service, typ, expected, actual string) Finding {
	return Finding{
		Service:  service,
		Type:     DriftType(typ),
		Expected: expected,
		Actual:   actual,
		Message:  typ + " drift",
		Severity: SeverityHigh,
	}
}

func TestExportReport_Empty(t *testing.T) {
	records := ExportReport(Report{}, exportTime)
	if len(records) != 0 {
		t.Fatalf("expected 0 records, got %d", len(records))
	}
}

func TestExportReport_PopulatesFields(t *testing.T) {
	r := Report{
		Findings: []Finding{
			exportFinding("web", "image", "nginx:1.25", "nginx:1.24"),
		},
	}
	records := ExportReport(r, exportTime)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	rec := records[0]
	if rec.Service != "web" {
		t.Errorf("service: got %q", rec.Service)
	}
	if rec.Timestamp != "2024-06-01T12:00:00Z" {
		t.Errorf("timestamp: got %q", rec.Timestamp)
	}
	if rec.Severity != "high" {
		t.Errorf("severity: got %q", rec.Severity)
	}
}

func TestExportReport_SortedByServiceThenType(t *testing.T) {
	r := Report{
		Findings: []Finding{
			exportFinding("web", "env", "A=1", "A=2"),
			exportFinding("api", "image", "x", "y"),
			exportFinding("web", "image", "x", "y"),
		},
	}
	records := ExportReport(r, exportTime)
	if records[0].Service != "api" {
		t.Errorf("first record should be api, got %q", records[0].Service)
	}
	if records[1].Type != "env" {
		t.Errorf("second record type should be env, got %q", records[1].Type)
	}
}

func TestWriteExport_JSON(t *testing.T) {
	records := []ExportRecord{{Service: "svc", Type: "image", Severity: "high"}}
	var buf bytes.Buffer
	if err := WriteExport(&buf, records, ExportJSON); err != nil {
		t.Fatal(err)
	}
	var out []ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 1 || out[0].Service != "svc" {
		t.Errorf("unexpected JSON output: %+v", out)
	}
}

func TestWriteExport_CSV(t *testing.T) {
	records := []ExportRecord{{Service: "svc", Type: "env", Severity: "low", Message: "env drift"}}
	var buf bytes.Buffer
	if err := WriteExport(&buf, records, ExportCSV); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected header + 1 row, got %d lines", len(lines))
	}
	if !strings.HasPrefix(lines[0], "timestamp") {
		t.Errorf("missing CSV header: %q", lines[0])
	}
}

func TestWriteExport_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := WriteExport(&buf, nil, ExportFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}
