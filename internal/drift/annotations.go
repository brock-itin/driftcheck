package drift

import (
	"fmt"
	"strings"
)

// Annotation holds a key-value metadata tag attached to a Finding.
type Annotation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AnnotateFindings attaches a set of annotations to every finding in the report
// whose service name matches the provided service filter. Pass an empty
// serviceName to annotate all findings.
func AnnotateFindings(r *Report, serviceName string, annotations []Annotation) *Report {
	if r == nil || len(annotations) == 0 {
		return r
	}

	updated := make([]Finding, 0, len(r.Findings))
	for _, f := range r.Findings {
		if serviceName == "" || f.Service == serviceName {
			f.Message = appendAnnotations(f.Message, annotations)
		}
		updated = append(updated, f)
	}

	return &Report{
		Findings: updated,
	}
}

// ParseAnnotations parses a slice of "key=value" strings into Annotation
// structs. Entries that do not contain "=" are skipped.
func ParseAnnotations(raw []string) ([]Annotation, error) {
	var annotations []Annotation
	for _, s := range raw {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid annotation %q: expected key=value format", s)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("invalid annotation %q: key must not be empty", s)
		}
		annotations = append(annotations, Annotation{Key: key, Value: val})
	}
	return annotations, nil
}

// appendAnnotations appends formatted annotations to an existing message string.
func appendAnnotations(msg string, annotations []Annotation) string {
	if len(annotations) == 0 {
		return msg
	}
	parts := make([]string, 0, len(annotations))
	for _, a := range annotations {
		parts = append(parts, fmt.Sprintf("%s=%s", a.Key, a.Value))
	}
	return fmt.Sprintf("%s [%s]", msg, strings.Join(parts, ", "))
}
