package drift

import "sort"

// FindingType classifies what kind of drift was detected.
type FindingType string

const (
	FindingTypeImage   FindingType = "image"
	FindingTypeEnv     FindingType = "env"
	FindingTypePorts   FindingType = "ports"
	FindingTypeMissing FindingType = "missing"
)

// Finding represents a single detected drift between a running container
// and its compose/helm definition.
type Finding struct {
	Service  string
	Type     FindingType
	Expected string
	Actual   string
	Severity Severity
}

// IsZero returns true if the Finding is the zero value.
func (f Finding) IsZero() bool {
	return f.Service == "" && f.Type == "" && f.Expected == "" && f.Actual == ""
}

// FilterBySeverity returns a new slice containing only findings at or above
// the given minimum severity level.
func FilterBySeverity(findings []Finding, min Severity) []Finding {
	result := make([]Finding, 0, len(findings))
	for _, f := range findings {
		if f.Severity >= min {
			result = append(result, f)
		}
	}
	return result
}

// BySeverity implements sort.Interface for []Finding, ordering by
// descending severity (Critical first).
type BySeverity []Finding

func (b BySeverity) Len() int      { return len(b) }
func (b BySeverity) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b BySeverity) Less(i, j int) bool {
	return b[i].Severity > b[j].Severity
}

// SortFindings sorts a slice of findings by descending severity in place.
func SortFindings(findings []Finding) {
	sort.Stable(BySeverity(findings))
}
