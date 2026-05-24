package output

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
)

func makeEnrichedFindings(n int) []drift.EnrichedFinding {
	out := make([]drift.EnrichedFinding, n)
	for i := range out {
		out[i] = drift.EnrichedFinding{
			Finding: drift.Finding{
				Service:   "svc",
				Type:      "image",
				Severity:  drift.SeverityHigh,
				Timestamp: time.Now().UTC(),
			},
			Environment: "prod",
			Cluster:     "k8s-us",
			RunbookURL:  "https://wiki.example.com/runbooks/image",
			EnrichedAt:  time.Now().UTC(),
		}
	}
	return out
}

func TestWriteEnriched_NoFindings(t *testing.T) {
	var sb strings.Builder
	if err := WriteEnriched(&sb, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "no enriched findings") {
		t.Errorf("expected no-findings message, got %q", sb.String())
	}
}

func TestWriteEnriched_WithFindings(t *testing.T) {
	findings := makeEnrichedFindings(2)
	var sb strings.Builder
	if err := WriteEnriched(&sb, findings); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	for _, want := range []string{"SERVICE", "ENVIRONMENT", "CLUSTER", "RUNBOOK", "prod", "k8s-us"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestEnrichmentStatusLine_NoFindings(t *testing.T) {
	line := EnrichmentStatusLine(nil, drift.EnrichmentSource{})
	if !strings.Contains(line, "no findings") {
		t.Errorf("unexpected status line: %q", line)
	}
}

func TestEnrichmentStatusLine_WithFindings(t *testing.T) {
	findings := makeEnrichedFindings(3)
	src := drift.EnrichmentSource{Environment: "staging", Cluster: "k8s-eu"}
	line := EnrichmentStatusLine(findings, src)
	if !strings.Contains(line, "3 finding(s)") {
		t.Errorf("expected count in line, got %q", line)
	}
	if !strings.Contains(line, "staging") {
		t.Errorf("expected env in line, got %q", line)
	}
	if !strings.Contains(line, "k8s-eu") {
		t.Errorf("expected cluster in line, got %q", line)
	}
}

func TestEnrichmentStatusLine_UnsetFields(t *testing.T) {
	findings := makeEnrichedFindings(1)
	src := drift.EnrichmentSource{}
	line := EnrichmentStatusLine(findings, src)
	if !strings.Contains(line, "(unset)") {
		t.Errorf("expected (unset) for empty env, got %q", line)
	}
}
