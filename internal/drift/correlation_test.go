package drift

import (
	"testing"
	"time"
)

func corrFinding(service, driftType string, at time.Time) Finding {
	return Finding{
		Service:    service,
		Type:       driftType,
		DetectedAt: at,
	}
}

func TestDefaultCorrelationOptions(t *testing.T) {
	opts := DefaultCorrelationOptions()
	if opts.Window != 5*time.Minute {
		t.Errorf("expected 5m window, got %v", opts.Window)
	}
	if opts.MinSize != 2 {
		t.Errorf("expected MinSize 2, got %d", opts.MinSize)
	}
}

func TestCorrelateFindings_Empty(t *testing.T) {
	result := CorrelateFindings(nil, DefaultCorrelationOptions())
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestCorrelateFindings_ZeroWindow(t *testing.T) {
	now := time.Now()
	findings := []Finding{
		corrFinding("svc-a", "image", now),
		corrFinding("svc-b", "env", now.Add(time.Second)),
	}
	opts := CorrelationOptions{Window: 0, MinSize: 2}
	result := CorrelateFindings(findings, opts)
	if result != nil {
		t.Errorf("expected nil for zero window, got %v", result)
	}
}

func TestCorrelateFindings_SingleBucket(t *testing.T) {
	now := time.Now()
	findings := []Finding{
		corrFinding("svc-a", "image", now),
		corrFinding("svc-b", "env", now.Add(30*time.Second)),
		corrFinding("svc-c", "label", now.Add(90*time.Second)),
	}
	opts := DefaultCorrelationOptions()
	result := CorrelateFindings(findings, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 window, got %d", len(result))
	}
	if len(result[0].Findings) != 3 {
		t.Errorf("expected 3 findings in window, got %d", len(result[0].Findings))
	}
	if len(result[0].Services) != 3 {
		t.Errorf("expected 3 services, got %d", len(result[0].Services))
	}
}

func TestCorrelateFindings_BelowMinSize(t *testing.T) {
	now := time.Now()
	findings := []Finding{
		corrFinding("svc-a", "image", now),
	}
	opts := DefaultCorrelationOptions()
	result := CorrelateFindings(findings, opts)
	if result != nil {
		t.Errorf("expected nil when below MinSize, got %v", result)
	}
}

func TestCorrelateFindings_SortedServices(t *testing.T) {
	now := time.Now()
	findings := []Finding{
		corrFinding("zebra", "image", now),
		corrFinding("alpha", "env", now.Add(10*time.Second)),
	}
	opts := DefaultCorrelationOptions()
	result := CorrelateFindings(findings, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 window, got %d", len(result))
	}
	if result[0].Services[0] != "alpha" || result[0].Services[1] != "zebra" {
		t.Errorf("expected sorted services, got %v", result[0].Services)
	}
}

func TestCorrelateFindings_TypesDeduped(t *testing.T) {
	now := time.Now()
	findings := []Finding{
		corrFinding("svc-a", "image", now),
		corrFinding("svc-b", "image", now.Add(10*time.Second)),
	}
	opts := DefaultCorrelationOptions()
	result := CorrelateFindings(findings, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 window, got %d", len(result))
	}
	if len(result[0].Types) != 1 || result[0].Types[0] != "image" {
		t.Errorf("expected deduplicated types [image], got %v", result[0].Types)
	}
}
