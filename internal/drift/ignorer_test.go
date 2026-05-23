package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempIgnoreList(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "ignore.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp ignore list: %v", err)
	}
	return p
}

func igFinding(service, typ, field string) Finding {
	return Finding{Service: service, Type: DriftType(typ), Field: field}
}

func TestLoadIgnoreList_MissingFile(t *testing.T) {
	il, err := LoadIgnoreList("/nonexistent/path/ignore.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(il.Rules) != 0 {
		t.Errorf("expected empty rules, got %d", len(il.Rules))
	}
}

func TestLoadIgnoreList_Valid(t *testing.T) {
	content := `ignore:
  - service: web
    type: image
  - service: ""
    type: env
    key: DEBUG
`
	p := writeTempIgnoreList(t, content)
	il, err := LoadIgnoreList(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(il.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(il.Rules))
	}
}

func TestApplyIgnoreList_NoRules(t *testing.T) {
	il := &IgnoreList{}
	findings := []Finding{igFinding("web", "image", ""), igFinding("db", "env", "PORT")}
	result := il.Apply(findings)
	if len(result) != 2 {
		t.Errorf("expected 2 findings, got %d", len(result))
	}
}

func TestApplyIgnoreList_ByService(t *testing.T) {
	il := &IgnoreList{Rules: []IgnoreRule{{Service: "web"}}}
	findings := []Finding{igFinding("web", "image", ""), igFinding("db", "env", "PORT")}
	result := il.Apply(findings)
	if len(result) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(result))
	}
	if result[0].Service != "db" {
		t.Errorf("expected db finding, got %s", result[0].Service)
	}
}

func TestApplyIgnoreList_ByTypeAndKey(t *testing.T) {
	il := &IgnoreList{Rules: []IgnoreRule{{Type: "env", Key: "DEBUG"}}}
	findings := []Finding{
		igFinding("web", "env", "DEBUG"),
		igFinding("web", "env", "PORT"),
		igFinding("web", "image", ""),
	}
	result := il.Apply(findings)
	if len(result) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(result))
	}
}

func TestApplyIgnoreList_NilList(t *testing.T) {
	var il *IgnoreList
	findings := []Finding{igFinding("web", "image", "")}
	result := il.Apply(findings)
	if len(result) != 1 {
		t.Errorf("expected 1 finding, got %d", len(result))
	}
}
