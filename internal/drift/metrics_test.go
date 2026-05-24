package drift

import (
	"testing"
	"time"
)

func metricFinding(service, typ string, sev Severity) Finding {
	return Finding{
		Service:  service,
		Type:     DriftType(typ),
		Severity: sev,
	}
}

func TestCollectMetrics_EmptyReport(t *testing.T) {
	r := Report{}
	m := CollectMetrics(r, 50*time.Millisecond)
	if m.FindingsTotal != 0 {
		t.Errorf("expected 0 findings, got %d", m.FindingsTotal)
	}
	if m.ServicesTotal != 0 {
		t.Errorf("expected 0 services, got %d", m.ServicesTotal)
	}
	if m.DurationMs != 50 {
		t.Errorf("expected 50ms, got %d", m.DurationMs)
	}
}

func TestCollectMetrics_WithFindings(t *testing.T) {
	r := Report{
		Findings: []Finding{
			metricFinding("web", "image", SeverityHigh),
			metricFinding("web", "env", SeverityLow),
			metricFinding("db", "image", SeverityHigh),
		},
	}
	m := CollectMetrics(r, 100*time.Millisecond)
	if m.FindingsTotal != 3 {
		t.Errorf("expected 3 findings, got %d", m.FindingsTotal)
	}
	if m.ServicesTotal != 2 {
		t.Errorf("expected 2 services, got %d", m.ServicesTotal)
	}
	if m.ServicesDrifted != 2 {
		t.Errorf("expected 2 drifted services, got %d", m.ServicesDrifted)
	}
	if m.BySeverity["high"] != 2 {
		t.Errorf("expected 2 high, got %d", m.BySeverity["high"])
	}
	if m.ByType["image"] != 2 {
		t.Errorf("expected 2 image findings, got %d", m.ByType["image"])
	}
}

func TestTopDriftedServices_Order(t *testing.T) {
	r := Report{
		Findings: []Finding{
			metricFinding("web", "image", SeverityHigh),
			metricFinding("web", "env", SeverityLow),
			metricFinding("db", "image", SeverityHigh),
			metricFinding("cache", "env", SeverityLow),
			metricFinding("cache", "image", SeverityHigh),
			metricFinding("cache", "env", SeverityLow),
		},
	}
	top := TopDriftedServices(r, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
	if top[0] != "cache (3)" {
		t.Errorf("expected cache first, got %s", top[0])
	}
	if top[1] != "web (2)" {
		t.Errorf("expected web second, got %s", top[1])
	}
}

func TestTopDriftedServices_FewerThanN(t *testing.T) {
	r := Report{
		Findings: []Finding{
			metricFinding("web", "image", SeverityHigh),
		},
	}
	top := TopDriftedServices(r, 5)
	if len(top) != 1 {
		t.Errorf("expected 1, got %d", len(top))
	}
}
