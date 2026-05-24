package output

import (
	"strings"
	"testing"

	"github.com/user/driftcheck/internal/drift"
)

func makeScores() []drift.DriftScore {
	return []drift.DriftScore{
		{Service: "web", Score: 15.0, FindingCount: 3},
		{Service: "db", Score: 5.0, FindingCount: 1},
	}
}

func TestWriteScores_NoScores(t *testing.T) {
	var sb strings.Builder
	WriteScores(&sb, nil)
	if !strings.Contains(sb.String(), "no drift scores") {
		t.Errorf("expected empty message, got: %q", sb.String())
	}
}

func TestWriteScores_WithScores(t *testing.T) {
	var sb strings.Builder
	WriteScores(&sb, makeScores())
	out := sb.String()

	for _, want := range []string{"web", "15.0", "db", "5.0", "TOTAL", "20.0"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q\ngot:\n%s", want, out)
		}
	}
}

func TestWriteScores_HeaderPresent(t *testing.T) {
	var sb strings.Builder
	WriteScores(&sb, makeScores())
	out := sb.String()
	if !strings.Contains(out, "SERVICE") || !strings.Contains(out, "SCORE") || !strings.Contains(out, "FINDINGS") {
		t.Errorf("expected header row in output, got:\n%s", out)
	}
}

func TestScoreThresholdLine_NoThreshold(t *testing.T) {
	line := ScoreThresholdLine(makeScores(), 0)
	if !strings.Contains(line, "no threshold set") {
		t.Errorf("unexpected line: %q", line)
	}
}

func TestScoreThresholdLine_WithinThreshold(t *testing.T) {
	line := ScoreThresholdLine(makeScores(), 50.0)
	if !strings.Contains(line, "within threshold") {
		t.Errorf("expected within-threshold message, got: %q", line)
	}
}

func TestScoreThresholdLine_ExceedsThreshold(t *testing.T) {
	line := ScoreThresholdLine(makeScores(), 10.0)
	if !strings.Contains(line, "EXCEEDS") {
		t.Errorf("expected EXCEEDS message, got: %q", line)
	}
}
