package drift

import (
	"strings"
	"testing"
)

func TestReport_NoFindings(t *testing.T) {
	var sb strings.Builder
	Report(&sb, nil)
	out := sb.String()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", out)
	}
}

func TestReport_WithFindings(t *testing.T) {
	findings := []Finding{
		{
			Service:  "web",
			Type:     DriftImage,
			Expected: "nginx:1.25",
			Actual:   "nginx:1.24",
			Message:  "image mismatch for service \"web\"",
		},
	}
	var sb strings.Builder
	Report(&sb, findings)
	out := sb.String()

	if !strings.Contains(out, "1 finding(s)") {
		t.Errorf("expected finding count in output, got: %s", out)
	}
	if !strings.Contains(out, "nginx:1.25") {
		t.Errorf("expected expected image in output, got: %s", out)
	}
	if !strings.Contains(out, "nginx:1.24") {
		t.Errorf("expected actual image in output, got: %s", out)
	}
}

func TestSummary_NoFindings(t *testing.T) {
	s := Summary(nil)
	if s != "no drift detected" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithFindings(t *testing.T) {
	findings := []Finding{
		{Service: "web", Type: DriftImage},
		{Service: "web", Type: DriftEnv},
		{Service: "db", Type: DriftImage},
	}
	s := Summary(findings)
	if !strings.Contains(s, "3 finding(s)") {
		t.Errorf("expected 3 findings in summary, got: %s", s)
	}
	if !strings.Contains(s, "2 service(s)") {
		t.Errorf("expected 2 services in summary, got: %s", s)
	}
}
