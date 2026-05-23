package drift

import (
	"sort"
	"testing"
)

func TestFinding_IsZero(t *testing.T) {
	var empty Finding
	if !empty.IsZero() {
		t.Error("zero-value Finding should be IsZero() == true")
	}

	populated := Finding{Service: "web", Type: "image_drift", Field: "image"}
	if populated.IsZero() {
		t.Error("populated Finding should be IsZero() == false")
	}
}

func TestBySeverity_Sort(t *testing.T) {
	findings := []Finding{
		{Service: "a", Severity: SeverityLow},
		{Service: "b", Severity: SeverityHigh},
		{Service: "c", Severity: SeverityMedium},
	}

	sort.Sort(BySeverity(findings))

	expected := []Severity{SeverityHigh, SeverityMedium, SeverityLow}
	for i, f := range findings {
		if f.Severity != expected[i] {
			t.Errorf("position %d: got severity %v, want %v", i, f.Severity, expected[i])
		}
	}
}

func TestBySeverity_StableForEqual(t *testing.T) {
	findings := []Finding{
		{Service: "x", Severity: SeverityMedium},
		{Service: "y", Severity: SeverityMedium},
	}

	sort.Stable(BySeverity(findings))

	if findings[0].Service != "x" || findings[1].Service != "y" {
		t.Error("stable sort should preserve relative order for equal severities")
	}
}

func TestBySeverity_Len(t *testing.T) {
	b := BySeverity([]Finding{{}, {}, {}})
	if b.Len() != 3 {
		t.Errorf("Len() = %d, want 3", b.Len())
	}
}
