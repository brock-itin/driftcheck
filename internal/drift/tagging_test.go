package drift

import (
	"testing"
)

func tagFinding(service, typ, field string, sev Severity) Finding {
	return Finding{
		Service:  service,
		Type:     DriftType(typ),
		Field:    field,
		Severity: sev,
	}
}

func TestTagFindings_NoRules(t *testing.T) {
	r := Report{Findings: []Finding{tagFinding("web", "image", "image", SeverityHigh)}}
	tags := TagFindings(r, nil)
	if len(tags) != 0 {
		t.Errorf("expected no tags, got %d", len(tags))
	}
}

func TestTagFindings_MatchByService(t *testing.T) {
	r := Report{Findings: []Finding{
		tagFinding("web", "image", "image", SeverityHigh),
		tagFinding("db", "env", "DB_PASS", SeverityLow),
	}}
	rules := []TagRule{{Tag: "frontend", Service: "web"}}
	tags := TagFindings(r, rules)
	if len(tags) != 1 {
		t.Fatalf("expected 1 tagged finding, got %d", len(tags))
	}
	for _, ts := range tags {
		if len(ts) != 1 || ts[0] != "frontend" {
			t.Errorf("unexpected tags: %v", ts)
		}
	}
}

func TestTagFindings_MatchByType(t *testing.T) {
	r := Report{Findings: []Finding{
		tagFinding("web", "image", "image", SeverityHigh),
		tagFinding("api", "env", "PORT", SeverityLow),
	}}
	rules := []TagRule{{Tag: "env-drift", Type: "env"}}
	tags := TagFindings(r, rules)
	if len(tags) != 1 {
		t.Fatalf("expected 1 tagged finding, got %d", len(tags))
	}
}

func TestTagFindings_MultipleTagsSameFinding(t *testing.T) {
	r := Report{Findings: []Finding{
		tagFinding("web", "image", "image", SeverityHigh),
	}}
	rules := []TagRule{
		{Tag: "critical", Severity: "high"},
		{Tag: "image-related", Type: "image"},
	}
	tags := TagFindings(r, rules)
	if len(tags) != 1 {
		t.Fatalf("expected 1 finding entry, got %d", len(tags))
	}
	for _, ts := range tags {
		if len(ts) != 2 {
			t.Errorf("expected 2 tags, got %v", ts)
		}
	}
}

func TestTagFindings_DeduplicatesTags(t *testing.T) {
	r := Report{Findings: []Finding{
		tagFinding("web", "image", "image", SeverityHigh),
	}}
	rules := []TagRule{
		{Tag: "important", Service: "web"},
		{Tag: "important", Type: "image"},
	}
	tags := TagFindings(r, rules)
	for _, ts := range tags {
		if len(ts) != 1 {
			t.Errorf("expected deduplication to 1 tag, got %v", ts)
		}
	}
}

func TestTagSummaryLine_NoTags(t *testing.T) {
	line := TagSummaryLine(nil)
	if line != "no tags applied" {
		t.Errorf("unexpected: %s", line)
	}
}

func TestTagSummaryLine_WithTags(t *testing.T) {
	tags := map[string]TagSet{
		"a::image::image": {"critical", "image-related"},
		"b::env::PORT":    {"env-drift"},
	}
	line := TagSummaryLine(tags)
	if line == "" || line == "no tags applied" {
		t.Errorf("expected non-empty summary, got: %s", line)
	}
}
