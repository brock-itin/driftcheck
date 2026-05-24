package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempPolicy(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "policy.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempPolicy: %v", err)
	}
	return p
}

func TestLoadPolicy_MissingFile(t *testing.T) {
	p, err := LoadPolicy("/nonexistent/policy.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Rules) != 0 {
		t.Errorf("expected empty policy, got %d rules", len(p.Rules))
	}
}

func TestLoadPolicy_Valid(t *testing.T) {
	path := writeTempPolicy(t, `rules:
  - type: image
    action: fail
  - type: env
    service: web
    action: ignore
`)
	p, err := LoadPolicy(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(p.Rules))
	}
	if p.Rules[0].Action != PolicyActionFail {
		t.Errorf("expected fail, got %s", p.Rules[0].Action)
	}
}

func TestPolicy_Resolve_Default(t *testing.T) {
	p := &Policy{}
	f := Finding{Type: "image", Service: "api"}
	if got := p.Resolve(f); got != PolicyActionWarn {
		t.Errorf("expected warn, got %s", got)
	}
}

func TestPolicy_Resolve_FirstMatch(t *testing.T) {
	p := &Policy{
		Rules: []PolicyRule{
			{Type: "image", Action: PolicyActionFail},
			{Type: "image", Action: PolicyActionIgnore},
		},
	}
	f := Finding{Type: "image", Service: "api"}
	if got := p.Resolve(f); got != PolicyActionFail {
		t.Errorf("expected fail, got %s", got)
	}
}

func TestPolicy_Apply_RemovesIgnored(t *testing.T) {
	p := &Policy{
		Rules: []PolicyRule{
			{Type: "env", Action: PolicyActionIgnore},
		},
	}
	findings := []Finding{
		{Type: "image", Service: "api"},
		{Type: "env", Service: "web"},
	}
	out := p.Apply(findings)
	if len(out) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(out))
	}
	if out[0].Type != "image" {
		t.Errorf("unexpected finding type: %s", out[0].Type)
	}
}

func TestPolicy_Apply_WildcardType(t *testing.T) {
	p := &Policy{
		Rules: []PolicyRule{
			{Type: "*", Action: PolicyActionIgnore},
		},
	}
	findings := []Finding{
		{Type: "image", Service: "api"},
		{Type: "env", Service: "web"},
	}
	out := p.Apply(findings)
	if len(out) != 0 {
		t.Errorf("expected 0 findings, got %d", len(out))
	}
}
