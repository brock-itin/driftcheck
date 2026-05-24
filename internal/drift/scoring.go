package drift

import (
	"sort"
)

// DriftScore represents a weighted drift score for a service.
type DriftScore struct {
	Service    string
	Score      float64
	FindingCount int
}

// severityWeight maps severity levels to numeric weights.
var severityWeight = map[string]float64{
	"critical": 10.0,
	"high":     5.0,
	"medium":   2.0,
	"low":      1.0,
	"info":     0.5,
}

// ScoreFindings computes a weighted drift score per service from a Report.
// Higher scores indicate more severe or numerous drift findings.
func ScoreFindings(r Report) []DriftScore {
	agg := make(map[string]float64)
	counts := make(map[string]int)

	for _, f := range r.Findings {
		w, ok := severityWeight[string(f.Severity)]
		if !ok {
			w = 1.0
		}
		agg[f.Service] += w
		counts[f.Service]++
	}

	scores := make([]DriftScore, 0, len(agg))
	for svc, total := range agg {
		scores = append(scores, DriftScore{
			Service:      svc,
			Score:        total,
			FindingCount: counts[svc],
		})
	}

	sort.Slice(scores, func(i, j int) bool {
		if scores[i].Score != scores[j].Score {
			return scores[i].Score > scores[j].Score
		}
		return scores[i].Service < scores[j].Service
	})

	return scores
}

// TotalScore sums all service scores into a single aggregate value.
func TotalScore(scores []DriftScore) float64 {
	var total float64
	for _, s := range scores {
		total += s.Score
	}
	return total
}
