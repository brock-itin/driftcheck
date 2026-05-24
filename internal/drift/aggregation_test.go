package drift

import (
	"testing"
	"time"
)

func aggFinding(svc, t string) Finding {
	return Finding{Service: svc, Type: t}
}

func makeEntry(ts time.Time, findings ...Finding) ChangelogEntry {
	return ChangelogEntry{Timestamp: ts, Findings: findings}
}

func TestAggregateByWindow_Empty(t *testing.T) {
	result := AggregateByWindow(nil, time.Hour)
	if len(result.Windows) != 0 {
		t.Errorf("expected 0 windows, got %d", len(result.Windows))
	}
	if result.TotalFindings != 0 {
		t.Errorf("expected 0 total, got %d", result.TotalFindings)
	}
}

func TestAggregateByWindow_ZeroDuration(t *testing.T) {
	now := time.Now()
	entries := []ChangelogEntry{makeEntry(now, aggFinding("svc", "image"))}
	result := AggregateByWindow(entries, 0)
	if len(result.Windows) != 0 {
		t.Errorf("expected 0 windows for zero duration")
	}
}

func TestAggregateByWindow_SingleBucket(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []ChangelogEntry{
		makeEntry(base, aggFinding("a", "image")),
		makeEntry(base.Add(10*time.Minute), aggFinding("b", "env")),
	}
	result := AggregateByWindow(entries, time.Hour)
	if len(result.Windows) != 1 {
		t.Fatalf("expected 1 window, got %d", len(result.Windows))
	}
	if result.TotalFindings != 2 {
		t.Errorf("expected 2 findings, got %d", result.TotalFindings)
	}
}

func TestAggregateByWindow_MultipleBuckets(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []ChangelogEntry{
		makeEntry(base, aggFinding("a", "image")),
		makeEntry(base.Add(2*time.Hour), aggFinding("b", "env")),
		makeEntry(base.Add(3*time.Hour), aggFinding("c", "label")),
	}
	result := AggregateByWindow(entries, time.Hour)
	if len(result.Windows) != 3 {
		t.Fatalf("expected 3 windows, got %d", len(result.Windows))
	}
	if result.TotalFindings != 3 {
		t.Errorf("expected 3 total findings, got %d", result.TotalFindings)
	}
}

func TestAggregateByWindow_PeakWindow(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []ChangelogEntry{
		makeEntry(base, aggFinding("a", "image")),
		makeEntry(base.Add(2*time.Hour), aggFinding("b", "env"), aggFinding("c", "label")),
	}
	result := AggregateByWindow(entries, time.Hour)
	if result.PeakWindow == nil {
		t.Fatal("expected non-nil peak window")
	}
	if len(result.PeakWindow.Findings) != 2 {
		t.Errorf("expected peak with 2 findings, got %d", len(result.PeakWindow.Findings))
	}
}

func TestAggregateByWindow_WindowBounds(t *testing.T) {
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	entries := []ChangelogEntry{makeEntry(base, aggFinding("x", "env"))}
	result := AggregateByWindow(entries, time.Hour)
	if len(result.Windows) != 1 {
		t.Fatalf("expected 1 window")
	}
	w := result.Windows[0]
	if !w.End.Equal(w.Start.Add(time.Hour)) {
		t.Errorf("window end should be start + 1h, got start=%v end=%v", w.Start, w.End)
	}
}
