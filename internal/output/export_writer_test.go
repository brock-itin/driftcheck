package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
)

var exportNow = time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

func makeExportReport(n int) drift.Report {
	r := drift.Report{}
	for i := 0; i < n; i++ {
		r.Findings = append(r.Findings, drift.Finding{
			Service:  "svc",
			Type:     drift.DriftTypeImage,
			Expected: "img:1",
			Actual:   "img:2",
			Severity: drift.SeverityHigh,
		})
	}
	return r
}

func TestWriteExportSummary_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	WriteExportSummary(&buf, nil, drift.ExportJSON, "out.json")
	if !strings.Contains(buf.String(), "no findings") {
		t.Errorf("expected 'no findings' message, got: %q", buf.String())
	}
}

func TestWriteExportSummary_WithFindings(t *testing.T) {
	records := []drift.ExportRecord{{Service: "svc"}, {Service: "api"}}
	var buf bytes.Buffer
	WriteExportSummary(&buf, records, drift.ExportCSV, "out.csv")
	out := buf.String()
	if !strings.Contains(out, "2") {
		t.Errorf("expected count 2 in output: %q", out)
	}
	if !strings.Contains(out, "out.csv") {
		t.Errorf("expected dest in output: %q", out)
	}
}

func TestExportAndWrite_JSON(t *testing.T) {
	var payload, summary bytes.Buffer
	r := makeExportReport(2)
	err := ExportAndWrite(&payload, &summary, r, drift.ExportJSON, "result.json", exportNow)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(payload.String(), "\"service\"") {
		t.Errorf("expected JSON output, got: %q", payload.String())
	}
	if !strings.Contains(summary.String(), "2") {
		t.Errorf("expected summary to mention 2 findings: %q", summary.String())
	}
}

func TestExportAndWrite_CSV(t *testing.T) {
	var payload, summary bytes.Buffer
	r := makeExportReport(1)
	err := ExportAndWrite(&payload, &summary, r, drift.ExportCSV, "", exportNow)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(payload.String()), "\n")
	if len(lines) < 2 {
		t.Errorf("expected header + data rows, got: %d lines", len(lines))
	}
}

func TestExportAndWrite_BadFormat(t *testing.T) {
	var payload, summary bytes.Buffer
	err := ExportAndWrite(&payload, &summary, drift.Report{}, drift.ExportFormat("toml"), "", exportNow)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
