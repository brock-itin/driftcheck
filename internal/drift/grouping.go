package drift

import "sort"

// GroupByType organises findings from a Report into a map keyed by drift type.
func GroupByType(r Report) map[string][]Finding {
	out := make(map[string][]Finding)
	for _, f := range r.Findings {
		out[f.Type] = append(out[f.Type], f)
	}
	return out
}

// GroupByService organises findings from a Report into a map keyed by service name.
func GroupByService(r Report) map[string][]Finding {
	out := make(map[string][]Finding)
	for _, f := range r.Findings {
		out[f.Service] = append(out[f.Service], f)
	}
	return out
}

// GroupBySeverity organises findings from a Report into a map keyed by severity label.
func GroupBySeverity(r Report) map[string][]Finding {
	out := make(map[string][]Finding)
	for _, f := range r.Findings {
		label := f.Severity.String()
		out[label] = append(out[label], f)
	}
	return out
}

// SortedGroupKeys returns the keys of a grouped-findings map in sorted order.
func SortedGroupKeys(m map[string][]Finding) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
