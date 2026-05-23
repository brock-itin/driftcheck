package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTempSuppressions(t *testing.T, sl SuppressionList) string {
	t.Helper()
	data, err := json.Marshal(sl)
	if err != nil {
		t.Fatalf("marshal suppressions: %v", err)
	}
	p := filepath.Join(t.TempDir(), "suppressions.json")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatalf("write suppressions: %v", err)
	}
	return p
}

func suppFinding(service, typ string) Finding {
	return Finding{Service: service, Type: FindingType(typ)}
}

func TestLoadSuppressions_MissingFile(t *testing.T) {
	sl, err := LoadSuppressions("/nonexistent/path/suppressions.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sl.Rules) != 0 {
		t.Errorf("expected empty rules, got %d", len(sl.Rules))
	}
}

func TestLoadSuppressions_Valid(t *testing.T) {
	original := SuppressionList{
		Rules: []SuppressionRule{
			{Service: "web", Type: "image", Reason: "pinned externally"},
		},
	}
	path := writeTempSuppressions(t, original)

	loaded, err := LoadSuppressions(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(loaded.Rules))
	}
	if loaded.Rules[0].Service != "web" {
		t.Errorf("expected service 'web', got %q", loaded.Rules[0].Service)
	}
}

func TestApply_NoRules(t *testing.T) {
	sl := &SuppressionList{}
	findings := []Finding{suppFinding("web", "image"), suppFinding("db", "env")}
	result := sl.Apply(findings)
	if len(result) != 2 {
		t.Errorf("expected 2 findings, got %d", len(result))
	}
}

func TestApply_MatchByServiceAndType(t *testing.T) {
	sl := &SuppressionList{
		Rules: []SuppressionRule{{Service: "web", Type: "image"}},
	}
	findings := []Finding{
		suppFinding("web", "image"),
		suppFinding("web", "env"),
		suppFinding("db", "image"),
	}
	result := sl.Apply(findings)
	if len(result) != 2 {
		t.Errorf("expected 2 findings, got %d", len(result))
	}
}

func TestApply_MatchByServiceOnly(t *testing.T) {
	sl := &SuppressionList{
		Rules: []SuppressionRule{{Service: "web"}},
	}
	findings := []Finding{
		suppFinding("web", "image"),
		suppFinding("web", "env"),
		suppFinding("db", "image"),
	}
	result := sl.Apply(findings)
	if len(result) != 1 {
		t.Errorf("expected 1 finding, got %d", len(result))
	}
	if result[0].Service != "db" {
		t.Errorf("expected db finding to remain")
	}
}

func TestApply_CaseInsensitive(t *testing.T) {
	sl := &SuppressionList{
		Rules: []SuppressionRule{{Service: "WEB", Type: "IMAGE"}},
	}
	findings := []Finding{suppFinding("web", "image")}
	result := sl.Apply(findings)
	if len(result) != 0 {
		t.Errorf("expected 0 findings, got %d", len(result))
	}
}
