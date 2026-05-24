package drift

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// PolicyAction defines what to do when a rule matches.
type PolicyAction string

const (
	PolicyActionWarn PolicyAction = "warn"
	PolicyActionFail PolicyAction = "fail"
	PolicyActionIgnore PolicyAction = "ignore"
)

// PolicyRule maps a drift type (and optional service) to an action.
type PolicyRule struct {
	Type    string       `yaml:"type"`
	Service string       `yaml:"service,omitempty"`
	Action  PolicyAction `yaml:"action"`
}

// Policy holds a list of rules evaluated in order.
type Policy struct {
	Rules []PolicyRule `yaml:"rules"`
}

// LoadPolicy reads a policy YAML file from path.
// If the file does not exist, an empty Policy is returned.
func LoadPolicy(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Policy{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("policy: read %s: %w", path, err)
	}
	var p Policy
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("policy: parse %s: %w", path, err)
	}
	return &p, nil
}

// Resolve returns the PolicyAction for a given finding.
// The first matching rule wins; default is PolicyActionWarn.
func (p *Policy) Resolve(f Finding) PolicyAction {
	for _, r := range p.Rules {
		typeMatch := r.Type == "*" || r.Type == string(f.Type)
		serviceMatch := r.Service == "" || r.Service == f.Service
		if typeMatch && serviceMatch {
			return r.Action
		}
	}
	return PolicyActionWarn
}

// Apply filters or annotates findings according to policy rules.
// Findings with action "ignore" are removed; others are returned as-is.
func (p *Policy) Apply(findings []Finding) []Finding {
	out := make([]Finding, 0, len(findings))
	for _, f := range findings {
		if p.Resolve(f) != PolicyActionIgnore {
			out = append(out, f)
		}
	}
	return out
}
