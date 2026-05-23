package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func baselineFinding(service, typ, field string) Finding {
	return Finding{
		Service:  service,
		Type:     typ,
		Field:    field,
		Expected: "expected",
		Actual:   "actual",
	}
}

func TestSaveAndLoadBaseline_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	r := Report{
		Findings: []Finding{
			baselineFinding("web", "image", "image"),
			baselineFinding("db", "env", "DB_HOST"),
		},
	}

	if err := SaveBaseline(path, r); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	b, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}

	if len(b.Findings) != len(r.Findings) {
		t.Errorf("findings count: got %d, want %d", len(b.Findings), len(r.Findings))
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveBaseline_BadPath(t *testing.T) {
	err := SaveBaseline("/nonexistent/dir/baseline.json", Report{})
	if err == nil {
		t.Error("expected error for bad path")
	}
}

func TestDiffAgainstBaseline_NewFindings(t *testing.T) {
	base := Baseline{
		CreatedAt: time.Now(),
		Findings:  []Finding{baselineFinding("web", "image", "image")},
	}
	current := Report{
		Findings: []Finding{
			baselineFinding("web", "image", "image"),
			baselineFinding("db", "env", "DB_HOST"),
		},
	}

	newF, resolved := DiffAgainstBaseline(base, current)
	if len(newF) != 1 {
		t.Errorf("new findings: got %d, want 1", len(newF))
	}
	if len(resolved) != 0 {
		t.Errorf("resolved: got %d, want 0", len(resolved))
	}
}

func TestDiffAgainstBaseline_Resolved(t *testing.T) {
	base := Baseline{
		CreatedAt: time.Now(),
		Findings: []Finding{
			baselineFinding("web", "image", "image"),
			baselineFinding("db", "env", "DB_HOST"),
		},
	}
	current := Report{
		Findings: []Finding{baselineFinding("web", "image", "image")},
	}

	newF, resolved := DiffAgainstBaseline(base, current)
	if len(newF) != 0 {
		t.Errorf("new findings: got %d, want 0", len(newF))
	}
	if len(resolved) != 1 {
		t.Errorf("resolved: got %d, want 1", len(resolved))
	}
}

func TestDiffAgainstBaseline_Empty(t *testing.T) {
	base := Baseline{CreatedAt: time.Now()}
	current := Report{}
	newF, resolved := DiffAgainstBaseline(base, current)
	if len(newF) != 0 || len(resolved) != 0 {
		t.Errorf("expected no diff for empty inputs")
	}
}

var _ = os.DevNull // ensure os imported
