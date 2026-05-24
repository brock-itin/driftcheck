package drift

import (
	"fmt"

	"github.com/your-org/driftcheck/internal/drift"
)

// DeduplicateFindings removes duplicate findings from a report, keeping the
// highest-severity instance when duplicates exist for the same key.
func DeduplicateFindings(r Report) Report {
	type entry struct {
		idx      int
		finding  Finding
	}

	seen := make(map[string]entry)

	for i, f := range r.Findings {
		key := dedupKey(f)
		existing, ok := seen[key]
		if !ok {
			seen[key] = entry{idx: i, finding: f}
			continue
		}
		// Keep whichever has higher severity.
		if f.Severity > existing.finding.Severity {
			seen[key] = entry{idx: i, finding: f}
		}
	}

	deduped := make([]Finding, 0, len(seen))
	for _, e := range seen {
		deduped = append(deduped, e.finding)
	}

	// Preserve deterministic ordering by sorting.
	sortFindings(deduped)

	return Report{
		Findings:  deduped,
		Timestamp: r.Timestamp,
	}
}

// dedupKey returns a string key that uniquely identifies a finding by its
// service, type, and field — ignoring severity and description.
func dedupKey(f Finding) string {
	return fmt.Sprintf("%s|%s|%s", f.Service, f.Type, f.Field)
}

// sortFindings sorts findings deterministically by service, then type, then field.
func sortFindings(findings []Finding) {
	for i := 1; i < len(findings); i++ {
		for j := i; j > 0; j-- {
			a, b := findings[j-1], findings[j]
			if a.Service > b.Service ||
				(a.Service == b.Service && a.Type > b.Type) ||
				(a.Service == b.Service && a.Type == b.Type && a.Field > b.Field) {
				findings[j-1], findings[j] = findings[j], findings[j-1]
			} else {
				break
			}
		}
	}
}

// DuplicateCount returns the number of duplicate findings that were removed.
func DuplicateCount(original, deduped Report) int {
	return len(original.Findings) - len(deduped.Findings)
}
