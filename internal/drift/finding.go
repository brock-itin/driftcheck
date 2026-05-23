package drift

// Finding describes a single detected drift between a compose definition
// and a running container.
type Finding struct {
	// Service is the compose service name associated with this finding.
	Service string `json:"service"`
	// Type categorises the drift (e.g. "image_drift", "env_drift").
	Type string `json:"type"`
	// Field is the specific field that differs (e.g. "image", "ENV_VAR").
	Field string `json:"field"`
	// Expected is the value declared in the compose definition.
	Expected string `json:"expected"`
	// Actual is the value observed in the running container.
	Actual string `json:"actual"`
	// Severity is the importance level assigned to this finding.
	Severity Severity `json:"severity"`
}

// IsZero reports whether the finding is empty / uninitialized.
func (f Finding) IsZero() bool {
	return f.Service == "" && f.Type == "" && f.Field == ""
}

// BySeverity implements sort.Interface for []Finding ordered by Severity
// descending (highest first).
type BySeverity []Finding

func (b BySeverity) Len() int           { return len(b) }
func (b BySeverity) Less(i, j int) bool { return b[i].Severity > b[j].Severity }
func (b BySeverity) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
