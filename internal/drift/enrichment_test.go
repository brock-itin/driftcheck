package drift

import (
	"strings"
	"testing"
	"time"
)

func enrichFinding(service, driftType string, sev Severity) Finding {
	return Finding{
		Service:   service,
		Type:      driftType,
		Severity:  sev,
		Timestamp: time.Now().UTC(),
	}
}

func TestEnrichFindings_Empty(t *testing.T) {
	out := EnrichFindings(nil, EnrichmentSource{Environment: "prod"})
	if len(out) != 0 {
		t.Fatalf("expected 0 enriched findings, got %d", len(out))
	}
}

func TestEnrichFindings_SetsEnvironmentAndCluster(t *testing.T) {
	findings := []Finding{
		enrichFinding("web", "image", SeverityHigh),
		enrichFinding("db", "env", SeverityLow),
	}
	src := EnrichmentSource{Environment: "staging", Cluster: "k8s-eu"}
	out := EnrichFindings(findings, src)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	for _, ef := range out {
		if ef.Environment != "staging" {
			t.Errorf("expected environment=staging, got %q", ef.Environment)
		}
		if ef.Cluster != "k8s-eu" {
			t.Errorf("expected cluster=k8s-eu, got %q", ef.Cluster)
		}
		if ef.EnrichedAt.IsZero() {
			t.Error("expected EnrichedAt to be set")
		}
	}
}

func TestEnrichFindings_RunbookURL(t *testing.T) {
	findings := []Finding{enrichFinding("api", "image_drift", SeverityHigh)}
	src := EnrichmentSource{RunbookBase: "https://wiki.example.com/runbooks"}
	out := EnrichFindings(findings, src)
	want := "https://wiki.example.com/runbooks/image-drift"
	if out[0].RunbookURL != want {
		t.Errorf("expected runbook %q, got %q", want, out[0].RunbookURL)
	}
}

func TestEnrichFindings_NoRunbookWhenBaseEmpty(t *testing.T) {
	findings := []Finding{enrichFinding("api", "env", SeverityLow)}
	out := EnrichFindings(findings, EnrichmentSource{})
	if out[0].RunbookURL != "" {
		t.Errorf("expected empty runbook URL, got %q", out[0].RunbookURL)
	}
}

func TestEnrichmentSummary_NoFindings(t *testing.T) {
	s := EnrichmentSummary(nil)
	if s != "no findings to enrich" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestEnrichmentSummary_WithFindings(t *testing.T) {
	findings := []Finding{
		enrichFinding("web", "image", SeverityHigh),
		enrichFinding("db", "env", SeverityLow),
	}
	src := EnrichmentSource{Environment: "prod"}
	out := EnrichFindings(findings, src)
	s := EnrichmentSummary(out)
	if !strings.Contains(s, "2 finding(s)") {
		t.Errorf("expected count in summary, got %q", s)
	}
	if !strings.Contains(s, "prod") {
		t.Errorf("expected environment in summary, got %q", s)
	}
}
