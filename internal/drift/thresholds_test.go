package drift

import (
	"strings"
	"testing"
)

func threshFinding(svc, typ string, sev Severity) Finding {
	return Finding{
		Service:  svc,
		Type:     typ,
		Severity: sev,
		Expected: "a",
		Actual:   "b",
	}
}

func TestEvaluateThresholds_NoBreaches(t *testing.T) {
	r := Report{Findings: []Finding{
		threshFinding("web", "image", SeverityHigh),
	}}
	scores := []ServiceScore{{Service: "web", Score: 5}}
	th := Threshold{MaxFindings: 10, MaxScore: 20, MaxPerService: 5, OnBreach: ThresholdFail}

	res := EvaluateThresholds(th, r, scores)
	if res.Breached {
		t.Fatalf("expected no breach, got messages: %v", res.Messages)
	}
}

func TestEvaluateThresholds_MaxFindingsBreached(t *testing.T) {
	findings := make([]Finding, 5)
	for i := range findings {
		findings[i] = threshFinding("web", "image", SeverityLow)
	}
	r := Report{Findings: findings}
	th := Threshold{MaxFindings: 3, OnBreach: ThresholdFail}

	res := EvaluateThresholds(th, r, nil)
	if !res.Breached {
		t.Fatal("expected breach")
	}
	if res.Action != ThresholdFail {
		t.Errorf("expected fail action, got %q", res.Action)
	}
	if !strings.Contains(res.Messages[0], "max_findings") {
		t.Errorf("unexpected message: %s", res.Messages[0])
	}
}

func TestEvaluateThresholds_MaxScoreBreached(t *testing.T) {
	r := Report{}
	scores := []ServiceScore{
		{Service: "web", Score: 15},
		{Service: "db", Score: 10},
	}
	th := Threshold{MaxScore: 20, OnBreach: ThresholdWarn}

	res := EvaluateThresholds(th, r, scores)
	if !res.Breached {
		t.Fatal("expected breach")
	}
	if res.Action != ThresholdWarn {
		t.Errorf("expected warn, got %q", res.Action)
	}
}

func TestEvaluateThresholds_MaxPerServiceBreached(t *testing.T) {
	r := Report{Findings: []Finding{
		threshFinding("api", "image", SeverityHigh),
		threshFinding("api", "env", SeverityMedium),
		threshFinding("api", "label", SeverityLow),
	}}
	th := Threshold{MaxPerService: 2, OnBreach: ThresholdFail}

	res := EvaluateThresholds(th, r, nil)
	if !res.Breached {
		t.Fatal("expected breach")
	}
	if !strings.Contains(res.Messages[0], "api") {
		t.Errorf("expected service name in message: %s", res.Messages[0])
	}
}

func TestEvaluateThresholds_DefaultActionIsWarn(t *testing.T) {
	r := Report{Findings: []Finding{
		threshFinding("x", "image", SeverityHigh),
		threshFinding("x", "env", SeverityHigh),
	}}
	th := Threshold{MaxFindings: 1} // OnBreach intentionally empty

	res := EvaluateThresholds(th, r, nil)
	if !res.Breached {
		t.Fatal("expected breach")
	}
	if res.Action != ThresholdWarn {
		t.Errorf("expected default warn, got %q", res.Action)
	}
}

func TestBreachSummary_NoBreaches(t *testing.T) {
	res := ThresholdResult{Breached: false}
	got := BreachSummary(res)
	if got != "no threshold breaches" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestBreachSummary_WithBreaches(t *testing.T) {
	res := ThresholdResult{
		Breached: true,
		Action:   ThresholdFail,
		Messages: []string{"finding count 5 exceeds max_findings 3"},
	}
	got := BreachSummary(res)
	if !strings.HasPrefix(got, "[FAIL]") {
		t.Errorf("expected [FAIL] prefix, got %q", got)
	}
	if !strings.Contains(got, "max_findings") {
		t.Errorf("expected message content, got %q", got)
	}
}
