package drift

import (
	"testing"
)

func TestSeverity_String(t *testing.T) {
	cases := []struct {
		sev  Severity
		want string
	}{
		{SeverityLow, "low"},
		{SeverityMedium, "medium"},
		{SeverityHigh, "high"},
		{Severity(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.sev.String(); got != tc.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tc.sev, got, tc.want)
		}
	}
}

func TestSeverityFor(t *testing.T) {
	cases := []struct {
		typ  string
		want Severity
	}{
		{"image_drift", SeverityHigh},
		{"env_drift", SeverityMedium},
		{"missing_container", SeverityHigh},
		{"unknown_type", SeverityLow},
		{"", SeverityLow},
	}
	for _, tc := range cases {
		if got := severityFor(tc.typ); got != tc.want {
			t.Errorf("severityFor(%q) = %v, want %v", tc.typ, got, tc.want)
		}
	}
}

func TestAnnotateWithSeverity(t *testing.T) {
	input := Report{
		Findings: []Finding{
			{Service: "web", Type: "image_drift", Field: "image"},
			{Service: "db", Type: "env_drift", Field: "ENV_VAR"},
			{Service: "cache", Type: "missing_container", Field: ""},
		},
	}

	result := AnnotateWithSeverity(input)

	expected := []Severity{SeverityHigh, SeverityMedium, SeverityHigh}
	for i, f := range result.Findings {
		if f.Severity != expected[i] {
			t.Errorf("finding[%d] severity = %v, want %v", i, f.Severity, expected[i])
		}
	}
}

func TestAnnotateWithSeverity_PreservesFields(t *testing.T) {
	input := Report{
		Findings: []Finding{
			{Service: "api", Type: "env_drift", Field: "PORT", Expected: "8080", Actual: "9090"},
		},
	}

	result := AnnotateWithSeverity(input)
	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(result.Findings))
	}
	f := result.Findings[0]
	if f.Service != "api" || f.Field != "PORT" || f.Expected != "8080" || f.Actual != "9090" {
		t.Errorf("fields not preserved: %+v", f)
	}
}
