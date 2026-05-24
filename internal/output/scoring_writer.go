package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/driftcheck/internal/drift"
)

// WriteScores renders a drift score table to w.
func WriteScores(w io.Writer, scores []drift.DriftScore) {
	if len(scores) == 0 {
		fmt.Fprintln(w, "no drift scores — all services are clean")
		return
	}

	fmt.Fprintf(w, "%-30s %8s %8s\n", "SERVICE", "SCORE", "FINDINGS")
	fmt.Fprintln(w, strings.Repeat("-", 50))
	for _, s := range scores {
		fmt.Fprintf(w, "%-30s %8.1f %8d\n", s.Service, s.Score, s.FindingCount)
	}
	fmt.Fprintln(w, strings.Repeat("-", 50))
	fmt.Fprintf(w, "%-30s %8.1f\n", "TOTAL", drift.TotalScore(scores))
}

// ScoreThresholdLine returns a human-readable threshold evaluation string.
// It reports whether the total drift score exceeds the given threshold.
func ScoreThresholdLine(scores []drift.DriftScore, threshold float64) string {
	total := drift.TotalScore(scores)
	if threshold <= 0 {
		return fmt.Sprintf("total drift score: %.1f (no threshold set)", total)
	}
	if total > threshold {
		return fmt.Sprintf("total drift score: %.1f — EXCEEDS threshold %.1f", total, threshold)
	}
	return fmt.Sprintf("total drift score: %.1f — within threshold %.1f", total, threshold)
}
