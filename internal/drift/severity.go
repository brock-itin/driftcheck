package drift

// Severity represents the importance level of a drift finding.
type Severity int

const (
	// SeverityLow indicates a minor configuration difference.
	SeverityLow Severity = iota
	// SeverityMedium indicates a notable configuration difference.
	SeverityMedium
	// SeverityHigh indicates a critical configuration difference.
	SeverityHigh
)

// String returns the human-readable name of the severity level.
func (s Severity) String() string {
	switch s {
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	default:
		return "unknown"
	}
}

// severityFor returns the appropriate Severity for a given finding type.
func severityFor(findingType string) Severity {
	switch findingType {
	case "image_drift":
		return SeverityHigh
	case "env_drift":
		return SeverityMedium
	case "missing_container":
		return SeverityHigh
	default:
		return SeverityLow
	}
}

// AnnotateWithSeverity sets the Severity field on each finding in the report
// based on its Type, returning a new copy of the report.
func AnnotateWithSeverity(r Report) Report {
	annotated := make([]Finding, len(r.Findings))
	for i, f := range r.Findings {
		f.Severity = severityFor(f.Type)
		annotated[i] = f
	}
	return Report{Findings: annotated}
}
