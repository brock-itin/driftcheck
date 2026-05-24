package drift

import (
	"testing"
	"time"
)

func dedupFinding(service, typ, field string, sev Severity) Finding {
	return Finding{
		Service:     service,
		Type:        typ,
		Field:       field,
		Expected:    "expected",
		Actual:      "actual",
		Severity:    sev,
		Description: "test finding",
	}
}

func TestDeduplicateFindings_NoDuplicates(t *testing.T) {
	r := Report{
		Timestamp: time.Now(),
		Findings: []Finding{
			dedupFinding("svcA", "image", "image", SeverityHigh),
			dedupFinding("svcB", "env", "PORT", SeverityLow),
		},
	}
	out := DeduplicateFindings(r)
	if len(out.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(out.Findings))
	}
}

func TestDeduplicateFindings_RemovesDuplicates(t *testing.T) {
	r := Report{
		Timestamp: time.Now(),
		Findings: []Finding{
			dedupFinding("svcA", "image", "image", SeverityLow),
			dedupFinding("svcA", "image", "image", SeverityLow), // exact duplicate
		},
	}
	out := DeduplicateFindings(r)
	if len(out.Findings) != 1 {
		t.Fatalf("expected 1 finding after dedup, got %d", len(out.Findings))
	}
}

func TestDeduplicateFindings_KeepsHighestSeverity(t *testing.T) {
	r := Report{
		Timestamp: time.Now(),
		Findings: []Finding{
			dedupFinding("svcA", "env", "PORT", SeverityLow),
			dedupFinding("svcA", "env", "PORT", SeverityHigh),
		},
	}
	out := DeduplicateFindings(r)
	if len(out.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(out.Findings))
	}
	if out.Findings[0].Severity != SeverityHigh {
		t.Errorf("expected SeverityHigh, got %v", out.Findings[0].Severity)
	}
}

func TestDeduplicateFindings_PreservesTimestamp(t *testing.T) {
	now := time.Now()
	r := Report{Timestamp: now, Findings: []Finding{}}
	out := DeduplicateFindings(r)
	if !out.Timestamp.Equal(now) {
		t.Errorf("timestamp not preserved")
	}
}

func TestDuplicateCount(t *testing.T) {
	original := Report{
		Findings: []Finding{
			dedupFinding("svcA", "image", "image", SeverityLow),
			dedupFinding("svcA", "image", "image", SeverityHigh),
			dedupFinding("svcB", "env", "PORT", SeverityLow),
		},
	}
	deduped := DeduplicateFindings(original)
	count := DuplicateCount(original, deduped)
	if count != 1 {
		t.Errorf("expected 1 duplicate removed, got %d", count)
	}
}

func TestDedupKey_Uniqueness(t *testing.T) {
	f1 := dedupFinding("svcA", "image", "image", SeverityHigh)
	f2 := dedupFinding("svcA", "image", "image", SeverityLow)
	f3 := dedupFinding("svcB", "image", "image", SeverityHigh)

	if dedupKey(f1) != dedupKey(f2) {
		t.Error("expected same key for same service/type/field")
	}
	if dedupKey(f1) == dedupKey(f3) {
		t.Error("expected different key for different service")
	}
}
