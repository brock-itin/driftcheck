package drift

import (
	"fmt"
	"strings"
	"time"
)

// EnrichmentSource provides metadata to attach to findings.
type EnrichmentSource struct {
	Environment string
	Cluster     string
	RunbookBase string
}

// EnrichedFinding wraps a Finding with additional operational context.
type EnrichedFinding struct {
	Finding
	Environment string    `json:"environment,omitempty"`
	Cluster     string    `json:"cluster,omitempty"`
	RunbookURL  string    `json:"runbook_url,omitempty"`
	EnrichedAt  time.Time `json:"enriched_at"`
}

// EnrichFindings attaches environment, cluster, and runbook metadata to each finding.
func EnrichFindings(findings []Finding, src EnrichmentSource) []EnrichedFinding {
	now := time.Now().UTC()
	enriched := make([]EnrichedFinding, 0, len(findings))
	for _, f := range findings {
		ef := EnrichedFinding{
			Finding:     f,
			Environment: src.Environment,
			Cluster:     src.Cluster,
			EnrichedAt:  now,
		}
		if src.RunbookBase != "" {
			ef.RunbookURL = buildRunbookURL(src.RunbookBase, f.Type)
		}
		enriched = append(enriched, ef)
	}
	return enriched
}

// EnrichmentSummary returns a human-readable summary line for a set of enriched findings.
func EnrichmentSummary(findings []EnrichedFinding) string {
	if len(findings) == 0 {
		return "no findings to enrich"
	}
	envs := uniqueStrings(findings, func(e EnrichedFinding) string { return e.Environment })
	return fmt.Sprintf("%d finding(s) enriched; environments: %s", len(findings), strings.Join(envs, ", "))
}

func buildRunbookURL(base, driftType string) string {
	slug := strings.ToLower(strings.ReplaceAll(driftType, "_", "-"))
	return strings.TrimRight(base, "/") + "/" + slug
}

func uniqueStrings(findings []EnrichedFinding, key func(EnrichedFinding) string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, f := range findings {
		v := key(f)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return []string{"(none)"}
	}
	return out
}
