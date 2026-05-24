package drift

import (
	"fmt"
	"sort"
	"strings"
)

// TagRule defines a rule that applies a tag to findings matching given criteria.
type TagRule struct {
	Tag      string `yaml:"tag"`
	Service  string `yaml:"service,omitempty"`
	Type     string `yaml:"type,omitempty"`
	Severity string `yaml:"severity,omitempty"`
}

// TagSet holds a deduplicated, sorted list of tags for a finding.
type TagSet []string

// TagFindings applies tag rules to all findings in the report, returning a
// map of finding key -> TagSet for each matched finding.
func TagFindings(r Report, rules []TagRule) map[string]TagSet {
	result := make(map[string]TagSet)
	for _, f := range r.Findings {
		key := fmt.Sprintf("%s::%s::%s", f.Service, f.Type, f.Field)
		for _, rule := range rules {
			if matchesTagRule(f, rule) {
				result[key] = appendUnique(result[key], rule.Tag)
			}
		}
	}
	for k := range result {
		sort.Strings(result[k])
	}
	return result
}

// TagSummaryLine returns a human-readable summary of tagging results.
func TagSummaryLine(tags map[string]TagSet) string {
	if len(tags) == 0 {
		return "no tags applied"
	}
	total := 0
	for _, ts := range tags {
		total += len(ts)
	}
	return fmt.Sprintf("%d tag(s) applied across %d finding(s)", total, len(tags))
}

func matchesTagRule(f Finding, rule TagRule) bool {
	if rule.Service != "" && !strings.EqualFold(f.Service, rule.Service) {
		return false
	}
	if rule.Type != "" && !strings.EqualFold(string(f.Type), rule.Type) {
		return false
	}
	if rule.Severity != "" && !strings.EqualFold(f.Severity.String(), rule.Severity) {
		return false
	}
	return true
}

func appendUnique(ts TagSet, tag string) TagSet {
	for _, t := range ts {
		if t == tag {
			return ts
		}
	}
	return append(ts, tag)
}
