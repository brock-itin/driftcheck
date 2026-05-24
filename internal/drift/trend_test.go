package drift

import (
	"testing"
	"time"
)

func makeEntry(t time.Time, findings []Finding) ChangelogEntry {
	return ChangelogEntry{
		RecordedAt: t,
		Findings:   findings,
	}
}

func trendFinding(typ DriftType) Finding {
	return Finding{Service: "svc", Type: typ, Expected: "a", Actual: "b"}
}

func TestBuildTrend_Empty(t *testing.T) {
	trend := BuildTrend(nil)
	if len(trend.Points) != 0 {
		t.Fatalf("expected 0 points, got %d", len(trend.Points))
	}
}

func TestBuildTrend_OrderedByTime(t *testing.T) {
	now := time.Now()
	entries := []ChangelogEntry{
		makeEntry(now.Add(2*time.Hour), []Finding{trendFinding(DriftTypeImage)}),
		makeEntry(now, []Finding{trendFinding(DriftTypeImage), trendFinding(DriftTypeEnv)}),
		makeEntry(now.Add(time.Hour), []Finding{trendFinding(DriftTypeEnv)}),
	}
	trend := BuildTrend(entries)
	if len(trend.Points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(trend.Points))
	}
	if !trend.Points[0].Timestamp.Equal(now) {
		t.Errorf("expected first point at now, got %v", trend.Points[0].Timestamp)
	}
}

func TestBuildTrend_ByTypeCounts(t *testing.T) {
	now := time.Now()
	entry := makeEntry(now, []Finding{
		trendFinding(DriftTypeImage),
		trendFinding(DriftTypeImage),
		trendFinding(DriftTypeEnv),
	})
	trend := BuildTrend([]ChangelogEntry{entry})
	p := trend.Points[0]
	if p.Total != 3 {
		t.Errorf("expected total 3, got %d", p.Total)
	}
	if p.ByType[string(DriftTypeImage)] != 2 {
		t.Errorf("expected 2 image findings, got %d", p.ByType[string(DriftTypeImage)])
	}
	if p.ByType[string(DriftTypeEnv)] != 1 {
		t.Errorf("expected 1 env finding, got %d", p.ByType[string(DriftTypeEnv)])
	}
}

func TestTrend_Delta(t *testing.T) {
	now := time.Now()
	entries := []ChangelogEntry{
		makeEntry(now, []Finding{trendFinding(DriftTypeImage)}),
		makeEntry(now.Add(time.Hour), []Finding{trendFinding(DriftTypeImage), trendFinding(DriftTypeEnv), trendFinding(DriftTypeEnv)}),
	}
	trend := BuildTrend(entries)
	if trend.Delta() != 2 {
		t.Errorf("expected delta 2, got %d", trend.Delta())
	}
}

func TestTrend_Delta_SinglePoint(t *testing.T) {
	trend := BuildTrend([]ChangelogEntry{makeEntry(time.Now(), nil)})
	if trend.Delta() != 0 {
		t.Errorf("expected delta 0 for single point, got %d", trend.Delta())
	}
}

func TestTrend_Latest_Empty(t *testing.T) {
	trend := BuildTrend(nil)
	if !trend.Latest().Timestamp.IsZero() {
		t.Error("expected zero timestamp for empty trend")
	}
}
