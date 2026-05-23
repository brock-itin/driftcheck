package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewChangelogEntry_EmptyReport(t *testing.T) {
	entry := NewChangelogEntry("docker-compose.yml", Report{})
	if entry.Total != 0 {
		t.Fatalf("expected 0 total, got %d", entry.Total)
	}
	if entry.ComposFile != "docker-compose.yml" {
		t.Fatalf("unexpected compose file: %s", entry.ComposFile)
	}
	if entry.Timestamp.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}

func TestNewChangelogEntry_WithFindings(t *testing.T) {
	r := Report{
		Findings: []Finding{
			{Service: "web", Type: "image", Severity: SeverityHigh},
			{Service: "db", Type: "env", Severity: SeverityLow},
			{Service: "web", Type: "env", Severity: SeverityHigh},
		},
	}
	entry := NewChangelogEntry("compose.yml", r)
	if entry.Total != 3 {
		t.Fatalf("expected 3 total, got %d", entry.Total)
	}
	if entry.BySeverity["high"] != 2 {
		t.Errorf("expected 2 high, got %d", entry.BySeverity["high"])
	}
	if entry.ByType["env"] != 2 {
		t.Errorf("expected 2 env findings, got %d", entry.ByType["env"])
	}
}

func TestLoadChangelog_MissingFile(t *testing.T) {
	cl, err := LoadChangelog(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cl.Entries) != 0 {
		t.Fatalf("expected empty changelog")
	}
}

func TestAppendChangelog_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "changelog.json")
	entry := ChangelogEntry{
		Timestamp:  time.Now().UTC().Truncate(time.Second),
		ComposFile: "docker-compose.yml",
		Total:      2,
		BySeverity: map[string]int{"high": 2},
		ByType:     map[string]int{"image": 1, "env": 1},
	}
	if err := AppendChangelog(path, entry); err != nil {
		t.Fatalf("append error: %v", err)
	}
	if err := AppendChangelog(path, entry); err != nil {
		t.Fatalf("second append error: %v", err)
	}
	cl, err := LoadChangelog(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(cl.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(cl.Entries))
	}
	if cl.Entries[0].Total != 2 {
		t.Errorf("unexpected total in first entry")
	}
}

func TestAppendChangelog_BadPath(t *testing.T) {
	err := AppendChangelog("/nonexistent/dir/changelog.json", ChangelogEntry{})
	if err == nil {
		t.Fatal("expected error for bad path")
	}
}

func TestLoadChangelog_Corrupt(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "changelog*.json")
	_, _ = f.WriteString("{not valid json")
	_ = f.Close()
	_, err := LoadChangelog(f.Name())
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
}
