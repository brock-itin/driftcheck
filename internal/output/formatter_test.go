package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftcheck/internal/drift"
	"github.com/driftcheck/internal/output"
)

func makeReport(findings []drift.Finding) *drift.Report {
	return &drift.Report{Findings: findings}
}

func TestWriteText_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatText)
	if err := f.Write(makeReport(nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %q", buf.String())
	}
}

func TestWriteText_WithFindings(t *testing.T) {
	findings := []drift.Finding{
		{Service: "web", Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"},
	}
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatText)
	if err := f.Write(makeReport(findings)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"web", "image", "nginx:1.25", "nginx:1.24"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
	}
}

func TestWriteJSON_WithFindings(t *testing.T) {
	findings := []drift.Finding{
		{Service: "db", Field: "env.DB_PASS", Expected: "secret", Actual: "wrong"},
	}
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatJSON)
	if err := f.Write(makeReport(findings)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result drift.Report
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(result.Findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(result.Findings))
	}
}

func TestFormatDefault_IsText(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, "")
	if err := f.Write(makeReport(nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected text fallback, got: %q", buf.String())
	}
}
