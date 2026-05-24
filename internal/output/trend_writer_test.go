package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/you/driftcheck/internal/drift"
)

func makeTrend(counts []int) drift.Trend {
	now := time.Now()
	entries := make([]drift.ChangelogEntry, len(counts))
	for i, c := range counts {
		findings := make([]drift.Finding, c)
		for j := range findings {
			findings[j] = drift.Finding{Service: "svc", Type: drift.DriftTypeImage, Expected: "a", Actual: "b"}
		}
		entries[i] = drift.ChangelogEntry{
			RecordedAt: now.Add(time.Duration(i) * time.Hour),
			Findings:   findings,
		}
	}
	return drift.BuildTrend(entries)
}

func TestWriteTrend_NoData(t *testing.T) {
	var buf bytes.Buffer
	WriteTrend(&buf, drift.Trend{})
	if !strings.Contains(buf.String(), "No trend data") {
		t.Errorf("expected no-data message, got: %s", buf.String())
	}
}

func TestWriteTrend_WithPoints(t *testing.T) {
	var buf bytes.Buffer
	WriteTrend(&buf, makeTrend([]int{1, 3, 2}))
	out := buf.String()
	if !strings.Contains(out, "3 snapshots") {
		t.Errorf("expected snapshot count in output, got: %s", out)
	}
	if !strings.Contains(out, "total=") {
		t.Errorf("expected total= in output, got: %s", out)
	}
}

func TestWriteTrend_DeltaPositive(t *testing.T) {
	var buf bytes.Buffer
	WriteTrend(&buf, makeTrend([]int{1, 4}))
	out := buf.String()
	if !strings.Contains(out, "+3") {
		t.Errorf("expected positive delta in output, got: %s", out)
	}
}

func TestWriteTrend_DeltaNegative(t *testing.T) {
	var buf bytes.Buffer
	WriteTrend(&buf, makeTrend([]int{5, 2}))
	out := buf.String()
	if !strings.Contains(out, "-3") {
		t.Errorf("expected negative delta in output, got: %s", out)
	}
}

func TestTrendDeltaLine_NoData(t *testing.T) {
	line := TrendDeltaLine(drift.Trend{})
	if line != "no data" {
		t.Errorf("expected 'no data', got %q", line)
	}
}

func TestTrendDeltaLine_NoChange(t *testing.T) {
	line := TrendDeltaLine(makeTrend([]int{3, 3}))
	if !strings.Contains(line, "no change") {
		t.Errorf("expected 'no change', got %q", line)
	}
}

func TestTrendDeltaLine_Increase(t *testing.T) {
	line := TrendDeltaLine(makeTrend([]int{2, 5}))
	if !strings.Contains(line, "+3") {
		t.Errorf("expected +3, got %q", line)
	}
}
