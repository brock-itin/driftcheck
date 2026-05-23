package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/driftcheck/internal/drift"
)

func makeTableReport(findings []drift.Finding) drift.Report {
	return drift.Report{Findings: findings}
}

func TestWriteTable_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	err := WriteTable(&buf, makeTableReport(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %q", buf.String())
	}
}

func TestWriteTable_WithFindings(t *testing.T) {
	findings := []drift.Finding{
		{
			Service:  "web",
			Type:     "image",
			Field:    "image",
			Expected: "nginx:1.25",
			Actual:   "nginx:1.24",
			Severity: drift.SeverityHigh,
		},
		{
			Service:  "db",
			Type:     "env",
			Field:    "DB_PORT",
			Expected: "5432",
			Actual:   "5433",
			Severity: drift.SeverityMedium,
		},
	}

	var buf bytes.Buffer
	err := WriteTable(&buf, makeTableReport(findings))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"SERVICE", "TYPE", "FIELD", "EXPECTED", "ACTUAL", "SEVERITY"} {
		if !strings.Contains(out, want) {
			t.Errorf("header missing %q in output:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "web") {
		t.Errorf("expected service 'web' in output")
	}
	if !strings.Contains(out, "nginx:1.25") {
		t.Errorf("expected expected image in output")
	}
	if !strings.Contains(out, "DB_PORT") {
		t.Errorf("expected field 'DB_PORT' in output")
	}
}

func TestWriteTable_Truncation(t *testing.T) {
	long := strings.Repeat("x", 50)
	findings := []drift.Finding{
		{
			Service:  "svc",
			Type:     "env",
			Field:    "KEY",
			Expected: long,
			Actual:   long,
			Severity: drift.SeverityLow,
		},
	}

	var buf bytes.Buffer
	if err := WriteTable(&buf, makeTableReport(findings)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(buf.String(), long) {
		t.Errorf("expected long string to be truncated in output")
	}
	if !strings.Contains(buf.String(), "…") {
		t.Errorf("expected truncation indicator '…' in output")
	}
}

func TestTruncate(t *testing.T) {
	cases := []struct {
		input string
		max   int
		want  string
	}{
		{"hello", 10, "hello"},
		{"hello world", 8, "hello w…"},
		{"", 5, ""},
		{"abc", 3, "abc"},
	}
	for _, tc := range cases {
		got := truncate(tc.input, tc.max)
		if got != tc.want {
			t.Errorf("truncate(%q, %d) = %q; want %q", tc.input, tc.max, got, tc.want)
		}
	}
}
