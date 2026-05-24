package drift

import (
	"testing"
)

func scoreFinding(service, severity, driftType string) Finding {
	return Finding{
		Service:  service,
		Type:     driftType,
		Severity: Severity(severity),
	}
}

func TestScoreFindings_EmptyReport(t *testing.T) {
	r := Report{}
	scores := ScoreFindings(r)
	if len(scores) != 0 {
		t.Fatalf("expected 0 scores, got %d", len(scores))
	}
}

func TestScoreFindings_SingleService(t *testing.T) {
	r := Report{
		Findings: []Finding{
			scoreFinding("web", "high", "image"),
			scoreFinding("web", "low", "env"),
		},
	}
	scores := ScoreFindings(r)
	if len(scores) != 1 {
		t.Fatalf("expected 1 score, got %d", len(scores))
	}
	if scores[0].Service != "web" {
		t.Errorf("expected service 'web', got %q", scores[0].Service)
	}
	// high=5.0 + low=1.0 = 6.0
	if scores[0].Score != 6.0 {
		t.Errorf("expected score 6.0, got %f", scores[0].Score)
	}
	if scores[0].FindingCount != 2 {
		t.Errorf("expected FindingCount 2, got %d", scores[0].FindingCount)
	}
}

func TestScoreFindings_OrderedByScore(t *testing.T) {
	r := Report{
		Findings: []Finding{
			scoreFinding("db", "low", "env"),
			scoreFinding("web", "critical", "image"),
			scoreFinding("cache", "medium", "label"),
		},
	}
	scores := ScoreFindings(r)
	if len(scores) != 3 {
		t.Fatalf("expected 3 scores, got %d", len(scores))
	}
	if scores[0].Service != "web" {
		t.Errorf("expected first service 'web', got %q", scores[0].Service)
	}
	if scores[1].Service != "cache" {
		t.Errorf("expected second service 'cache', got %q", scores[1].Service)
	}
	if scores[2].Service != "db" {
		t.Errorf("expected third service 'db', got %q", scores[2].Service)
	}
}

func TestScoreFindings_UnknownSeverityDefaultsToOne(t *testing.T) {
	r := Report{
		Findings: []Finding{
			{Service: "svc", Severity: Severity("unknown"), Type: "env"},
		},
	}
	scores := ScoreFindings(r)
	if len(scores) != 1 {
		t.Fatalf("expected 1 score, got %d", len(scores))
	}
	if scores[0].Score != 1.0 {
		t.Errorf("expected score 1.0 for unknown severity, got %f", scores[0].Score)
	}
}

func TestTotalScore_Empty(t *testing.T) {
	if TotalScore(nil) != 0 {
		t.Error("expected 0 for nil scores")
	}
}

func TestTotalScore_Sums(t *testing.T) {
	scores := []DriftScore{
		{Service: "a", Score: 5.0},
		{Service: "b", Score: 3.0},
		{Service: "c", Score: 2.0},
	}
	if TotalScore(scores) != 10.0 {
		t.Errorf("expected total 10.0, got %f", TotalScore(scores))
	}
}
