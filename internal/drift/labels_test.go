package drift

import (
	"testing"

	"github.com/yourorg/driftcheck/internal/compose"
)

func labelService(name string, labels map[string]string) compose.Service {
	return compose.Service{Name: name, Labels: labels}
}

func TestDetectLabelDrift_NoLabels(t *testing.T) {
	svc := labelService("web", map[string]string{})
	got := DetectLabelDrift(svc, map[string]string{"app": "web"})
	if len(got) != 0 {
		t.Fatalf("expected no drift, got %d", len(got))
	}
}

func TestDetectLabelDrift_Match(t *testing.T) {
	svc := labelService("web", map[string]string{"env": "prod"})
	got := DetectLabelDrift(svc, map[string]string{"env": "prod"})
	if len(got) != 0 {
		t.Fatalf("expected no drift, got %d", len(got))
	}
}

func TestDetectLabelDrift_ValueMismatch(t *testing.T) {
	svc := labelService("web", map[string]string{"env": "prod"})
	got := DetectLabelDrift(svc, map[string]string{"env": "staging"})
	if len(got) != 1 {
		t.Fatalf("expected 1 drift, got %d", len(got))
	}
	if got[0].Label != "env" {
		t.Errorf("unexpected label %q", got[0].Label)
	}
	if got[0].Expected != "prod" || got[0].Actual != "staging" {
		t.Errorf("unexpected values: expected=%q actual=%q", got[0].Expected, got[0].Actual)
	}
}

func TestDetectLabelDrift_MissingLabel(t *testing.T) {
	svc := labelService("api", map[string]string{"version": "1.2"})
	got := DetectLabelDrift(svc, map[string]string{})
	if len(got) != 1 {
		t.Fatalf("expected 1 drift, got %d", len(got))
	}
	if !got[0].Missing {
		t.Error("expected Missing=true")
	}
}

func TestDetectLabelDrift_CaseInsensitive(t *testing.T) {
	svc := labelService("svc", map[string]string{"tier": "Frontend"})
	got := DetectLabelDrift(svc, map[string]string{"tier": "frontend"})
	if len(got) != 0 {
		t.Fatalf("expected no drift (case-insensitive), got %d", len(got))
	}
}

func TestLabelDriftToFindings_Empty(t *testing.T) {
	findings := LabelDriftToFindings(nil)
	if len(findings) != 0 {
		t.Fatalf("expected empty findings, got %d", len(findings))
	}
}

func TestLabelDriftToFindings_MapsCorrectly(t *testing.T) {
	drifts := []LabelDrift{
		{Service: "web", Label: "env", Expected: "prod", Actual: "staging"},
		{Service: "db", Label: "owner", Expected: "team-a", Missing: true},
	}
	findings := LabelDriftToFindings(drifts)
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
	if findings[0].Type != "label" {
		t.Errorf("expected type 'label', got %q", findings[0].Type)
	}
	if findings[1].Message != "label missing on container" {
		t.Errorf("unexpected message: %q", findings[1].Message)
	}
}
