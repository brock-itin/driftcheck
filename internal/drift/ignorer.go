package drift

import (
	"strings"

	"gopkg.in/yaml.v3"
	"os"
)

// IgnoreRule defines a single ignore rule for drift findings.
type IgnoreRule struct {
	// Service is the service name to ignore. Empty means all services.
	Service string `yaml:"service"`
	// Type is the drift type to ignore (e.g. "image", "env"). Empty means all types.
	Type string `yaml:"type"`
	// Key is the specific env var or field key to ignore. Empty means all keys.
	Key string `yaml:"key"`
}

// IgnoreList holds a collection of ignore rules.
type IgnoreList struct {
	Rules []IgnoreRule `yaml:"ignore"`
}

// LoadIgnoreList reads an ignore list from a YAML file at path.
// If the file does not exist, an empty IgnoreList is returned.
func LoadIgnoreList(path string) (*IgnoreList, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &IgnoreList{}, nil
		}
		return nil, err
	}
	var il IgnoreList
	if err := yaml.Unmarshal(data, &il); err != nil {
		return nil, err
	}
	return &il, nil
}

// Apply filters out findings that match any rule in the IgnoreList.
func (il *IgnoreList) Apply(findings []Finding) []Finding {
	if il == nil || len(il.Rules) == 0 {
		return findings
	}
	out := findings[:0:0]
	for _, f := range findings {
		if !il.matches(f) {
			out = append(out, f)
		}
	}
	return out
}

// matches returns true if any rule in the list matches the finding.
func (il *IgnoreList) matches(f Finding) bool {
	for _, r := range il.Rules {
		if r.matches(f) {
			return true
		}
	}
	return false
}

// matches returns true if the rule applies to the given finding.
func (r IgnoreRule) matches(f Finding) bool {
	if r.Service != "" && !strings.EqualFold(r.Service, f.Service) {
		return false
	}
	if r.Type != "" && !strings.EqualFold(r.Type, string(f.Type)) {
		return false
	}
	if r.Key != "" && !strings.EqualFold(r.Key, f.Field) {
		return false
	}
	return true
}
