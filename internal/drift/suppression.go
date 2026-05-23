package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// SuppressionRule defines a rule that silences specific drift findings.
type SuppressionRule struct {
	Service string `json:"service"` // empty means match all services
	Type    string `json:"type"`    // empty means match all types
	Reason  string `json:"reason"`  // optional human-readable justification
}

// SuppressionList holds a collection of suppression rules.
type SuppressionList struct {
	Rules []SuppressionRule `json:"suppressions"`
}

// LoadSuppressions reads a suppression file from disk.
// Returns an empty list if the file does not exist.
func LoadSuppressions(path string) (*SuppressionList, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &SuppressionList{}, nil
		}
		return nil, fmt.Errorf("read suppression file: %w", err)
	}

	var sl SuppressionList
	if err := json.Unmarshal(data, &sl); err != nil {
		return nil, fmt.Errorf("parse suppression file: %w", err)
	}
	return &sl, nil
}

// Apply removes findings that match any suppression rule.
func (sl *SuppressionList) Apply(findings []Finding) []Finding {
	if len(sl.Rules) == 0 {
		return findings
	}

	out := findings[:0:0]
	for _, f := range findings {
		if !sl.suppresses(f) {
			out = append(out, f)
		}
	}
	return out
}

// suppresses returns true when at least one rule matches the finding.
func (sl *SuppressionList) suppresses(f Finding) bool {
	for _, r := range sl.Rules {
		serviceMatch := r.Service == "" || strings.EqualFold(r.Service, f.Service)
		typeMatch := r.Type == "" || strings.EqualFold(r.Type, string(f.Type))
		if serviceMatch && typeMatch {
			return true
		}
	}
	return false
}
