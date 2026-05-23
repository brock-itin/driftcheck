package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ChangelogEntry records a drift detection run and its findings summary.
type ChangelogEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	ComposFile string    `json:"compose_file"`
	Total      int       `json:"total_findings"`
	BySeverity map[string]int `json:"by_severity"`
	ByType     map[string]int `json:"by_type"`
}

// Changelog holds an ordered list of detection run entries.
type Changelog struct {
	Entries []ChangelogEntry `json:"entries"`
}

// AppendChangelog loads an existing changelog (if any), appends the new entry, and saves it.
func AppendChangelog(path string, entry ChangelogEntry) error {
	cl, err := LoadChangelog(path)
	if err != nil {
		return fmt.Errorf("changelog load: %w", err)
	}
	cl.Entries = append(cl.Entries, entry)
	data, err := json.MarshalIndent(cl, "", "  ")
	if err != nil {
		return fmt.Errorf("changelog marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("changelog write: %w", err)
	}
	return nil
}

// LoadChangelog reads a changelog file; returns an empty Changelog if the file does not exist.
func LoadChangelog(path string) (*Changelog, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Changelog{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("changelog read: %w", err)
	}
	var cl Changelog
	if err := json.Unmarshal(data, &cl); err != nil {
		return nil, fmt.Errorf("changelog unmarshal: %w", err)
	}
	return &cl, nil
}

// NewChangelogEntry builds a ChangelogEntry from a Report.
func NewChangelogEntry(composeFile string, r Report) ChangelogEntry {
	bySeverity := make(map[string]int)
	byType := make(map[string]int)
	for _, f := range r.Findings {
		bySeverity[f.Severity.String()]++
		byType[f.Type]++
	}
	return ChangelogEntry{
		Timestamp:  time.Now().UTC(),
		ComposFile: composeFile,
		Total:      len(r.Findings),
		BySeverity: bySeverity,
		ByType:     byType,
	}
}
