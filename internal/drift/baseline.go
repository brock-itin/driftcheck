package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved drift report used for comparison across runs.
type Baseline struct {
	CreatedAt time.Time `json:"created_at"`
	Findings  []Finding `json:"findings"`
}

// SaveBaseline writes the findings from a Report to a JSON file at path.
func SaveBaseline(path string, r Report) error {
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Findings:  r.Findings,
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("baseline: create %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(b); err != nil {
		return fmt.Errorf("baseline: encode: %w", err)
	}
	return nil
}

// LoadBaseline reads a Baseline from a JSON file at path.
func LoadBaseline(path string) (Baseline, error) {
	var b Baseline
	f, err := os.Open(path)
	if err != nil {
		return b, fmt.Errorf("baseline: open %q: %w", path, err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return b, fmt.Errorf("baseline: decode: %w", err)
	}
	return b, nil
}

// DiffAgainstBaseline returns findings present in current that were not in base,
// and findings that were in base but are no longer present (resolved).
func DiffAgainstBaseline(base Baseline, current Report) (newFindings, resolved []Finding) {
	baseSet := make(map[string]struct{}, len(base.Findings))
	for _, f := range base.Findings {
		baseSet[findingKey(f)] = struct{}{}
	}

	currentSet := make(map[string]struct{}, len(current.Findings))
	for _, f := range current.Findings {
		currentSet[findingKey(f)] = struct{}{}
		if _, ok := baseSet[findingKey(f)]; !ok {
			newFindings = append(newFindings, f)
		}
	}

	for _, f := range base.Findings {
		if _, ok := currentSet[findingKey(f)]; !ok {
			resolved = append(resolved, f)
		}
	}
	return newFindings, resolved
}

func findingKey(f Finding) string {
	return fmt.Sprintf("%s|%s|%s", f.Service, f.Type, f.Field)
}
