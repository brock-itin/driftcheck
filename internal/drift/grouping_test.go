package drift

import (
	"testing"
)

func groupFinding(service, driftType string, sev Severity) Finding {
	return Finding{
		Service:  service,
		Type:     driftType,
		Severity: sev,
		Field:    "test",
		Want:     "a",
		Got:      "b",
	}
}

func TestGroupByType_Empty(t *testing.T) {
	r := Report{}
	got := GroupByType(r)
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %d keys", len(got))
	}
}

func TestGroupByType_Groups(t *testing.T) {
	r := Report{
		Findings: []Finding{
			groupFinding("web", "image", SeverityHigh),
			groupFinding("db", "image", SeverityHigh),
			groupFinding("web", "env", SeverityLow),
		},
	}
	got := GroupByType(r)
	if len(got["image"]) != 2 {
		t.Errorf("expected 2 image findings, got %d", len(got["image"]))
	}
	if len(got["env"]) != 1 {
		t.Errorf("expected 1 env finding, got %d", len(got["env"]))
	}
}

func TestGroupByService_Groups(t *testing.T) {
	r := Report{
		Findings: []Finding{
			groupFinding("web", "image", SeverityHigh),
			groupFinding("web", "env", SeverityLow),
			groupFinding("db", "image", SeverityHigh),
		},
	}
	got := GroupByService(r)
	if len(got["web"]) != 2 {
		t.Errorf("expected 2 findings for web, got %d", len(got["web"]))
	}
	if len(got["db"]) != 1 {
		t.Errorf("expected 1 finding for db, got %d", len(got["db"]))
	}
}

func TestGroupBySeverity_Groups(t *testing.T) {
	r := Report{
		Findings: []Finding{
			groupFinding("web", "image", SeverityHigh),
			groupFinding("db", "env", SeverityLow),
			groupFinding("cache", "image", SeverityHigh),
		},
	}
	got := GroupBySeverity(r)
	if len(got[SeverityHigh.String()]) != 2 {
		t.Errorf("expected 2 high findings, got %d", len(got[SeverityHigh.String()]))
	}
	if len(got[SeverityLow.String()]) != 1 {
		t.Errorf("expected 1 low finding, got %d", len(got[SeverityLow.String()]))
	}
}

func TestSortedGroupKeys_Order(t *testing.T) {
	m := map[string][]Finding{
		"zebra": {},
		"alpha": {},
		"mango": {},
	}
	keys := SortedGroupKeys(m)
	expected := []string{"alpha", "mango", "zebra"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: want %q, got %q", i, expected[i], k)
		}
	}
}
