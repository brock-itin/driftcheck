package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftcheck/internal/drift"
)

func makeSummaryReport(findings []drift.Finding) drift.Report {
	return drift.Report{Findings: findings}
}

func TestWriteSummaryTable_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	WriteSummaryTable(&buf, makeSummaryReport(nil))
	out := buf.String()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", out)
	}
}

func TestWriteSummaryTable_WithFindings(t *testing.T) {
	findings := []drift.Finding{
		{Service: "web", Type: "image", Field: "image", Expected: "nginx:1.24", Actual: "nginx:1.23"},
		{Service: "web", Type: "env", Field: "PORT", Expected: "8080", Actual: "9090"},
		{Service: "db", Type: "image", Field: "image", Expected: "postgres:15", Actual: "postgres:14"},
	}
	var buf bytes.Buffer
	WriteSummaryTable(&buf, makeSummaryReport(findings))
	out := buf.String()

	if !strings.Contains(out, "Drift Summary") {
		t.Errorf("expected header, got: %s", out)
	}
	if !strings.Contains(out, "web") {
		t.Errorf("expected service 'web' in output")
	}
	if !strings.Contains(out, "db") {
		t.Errorf("expected service 'db' in output")
	}
	if !strings.Contains(out, "3 finding(s)") {
		t.Errorf("expected finding count in header, got: %s", out)
	}
	if !strings.Contains(out, "2 service(s)") {
		t.Errorf("expected service count in header, got: %s", out)
	}
}

func TestServiceDriftStatus_NoFindings(t *testing.T) {
	status := ServiceDriftStatus(makeSummaryReport(nil))
	if len(status) != 0 {
		t.Errorf("expected empty status map, got %v", status)
	}
}

func TestServiceDriftStatus_WithFindings(t *testing.T) {
	findings := []drift.Finding{
		{Service: "web", Type: "image"},
		{Service: "web", Type: "env"},
		{Service: "db", Type: "image"},
	}
	status := ServiceDriftStatus(makeSummaryReport(findings))
	if !status["web"] {
		t.Errorf("expected web to have drift")
	}
	if !status["db"] {
		t.Errorf("expected db to have drift")
	}
	if _, ok := status["cache"]; ok {
		t.Errorf("cache should not appear in status")
	}
}
