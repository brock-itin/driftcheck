package drift

import (
	"fmt"
	"strings"
)

// ThresholdAction defines what to do when a threshold is breached.
type ThresholdAction string

const (
	ThresholdWarn ThresholdAction = "warn"
	ThresholdFail ThresholdAction = "fail"
)

// Threshold defines a numeric limit for a drift metric.
type Threshold struct {
	MaxFindings  int            `yaml:"max_findings"`
	MaxScore     float64        `yaml:"max_score"`
	MaxPerService int           `yaml:"max_per_service"`
	OnBreach     ThresholdAction `yaml:"on_breach"`
}

// ThresholdResult holds the outcome of evaluating a threshold.
type ThresholdResult struct {
	Breached bool
	Action   ThresholdAction
	Messages []string
}

// EvaluateThresholds checks the given report and scores against the threshold config.
func EvaluateThresholds(t Threshold, r Report, scores []ServiceScore) ThresholdResult {
	var msgs []string

	totalFindings := len(r.Findings)
	if t.MaxFindings > 0 && totalFindings > t.MaxFindings {
		msgs = append(msgs, fmt.Sprintf(
			"finding count %d exceeds max_findings %d",
			totalFindings, t.MaxFindings,
		))
	}

	if t.MaxScore > 0 {
		total := TotalScore(scores)
		if total > t.MaxScore {
			msgs = append(msgs, fmt.Sprintf(
				"total drift score %.1f exceeds max_score %.1f",
				total, t.MaxScore,
			))
		}
	}

	if t.MaxPerService > 0 {
		byService := GroupByService(r)
		for svc, findings := range byService {
			if len(findings) > t.MaxPerService {
				msgs = append(msgs, fmt.Sprintf(
					"service %q has %d findings, exceeds max_per_service %d",
					svc, len(findings), t.MaxPerService,
				))
			}
		}
	}

	if len(msgs) == 0 {
		return ThresholdResult{Breached: false}
	}

	action := t.OnBreach
	if action == "" {
		action = ThresholdWarn
	}

	return ThresholdResult{
		Breached: true,
		Action:   action,
		Messages: msgs,
	}
}

// BreachSummary returns a human-readable summary of breached thresholds.
func BreachSummary(res ThresholdResult) string {
	if !res.Breached {
		return "no threshold breaches"
	}
	return fmt.Sprintf("[%s] %s", strings.ToUpper(string(res.Action)), strings.Join(res.Messages, "; "))
}
