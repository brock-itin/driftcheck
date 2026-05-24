package drift

import (
	"testing"
)

func normFinding(svc, typ, expected, actual string) Finding {
	return Finding{
		Service:  svc,
		Type:     typ,
		Expected: expected,
		Actual:   actual,
	}
}

func TestDefaultNormalizeOptions(t *testing.T) {
	opts := DefaultNormalizeOptions()
	if !opts.TrimEnvWhitespace {
		t.Error("expected TrimEnvWhitespace to be true")
	}
	if !opts.LowercaseImageRefs {
		t.Error("expected LowercaseImageRefs to be true")
	}
	if !opts.SortEnvKeys {
		t.Error("expected SortEnvKeys to be true")
	}
}

func TestNormalizeFindings_ImageLowercase(t *testing.T) {
	findings := []Finding{
		normFinding("web", "image", "Nginx:Latest", "NGINX:1.25"),
	}
	opts := DefaultNormalizeOptions()
	result := NormalizeFindings(findings, opts)
	if result[0].Expected != "nginx:latest" {
		t.Errorf("expected 'nginx:latest', got %q", result[0].Expected)
	}
	if result[0].Actual != "nginx:1.25" {
		t.Errorf("expected 'nginx:1.25', got %q", result[0].Actual)
	}
}

func TestNormalizeFindings_EnvTrim(t *testing.T) {
	findings := []Finding{
		normFinding("api", "env", "  production  ", "\tstaging\t"),
	}
	opts := DefaultNormalizeOptions()
	result := NormalizeFindings(findings, opts)
	if result[0].Expected != "production" {
		t.Errorf("expected 'production', got %q", result[0].Expected)
	}
	if result[0].Actual != "staging" {
		t.Errorf("expected 'staging', got %q", result[0].Actual)
	}
}

func TestNormalizeFindings_NoOpWhenDisabled(t *testing.T) {
	findings := []Finding{
		normFinding("web", "image", "Nginx:Latest", "NGINX:1.25"),
	}
	opts := NormalizeOptions{LowercaseImageRefs: false}
	result := NormalizeFindings(findings, opts)
	if result[0].Expected != "Nginx:Latest" {
		t.Errorf("expected original case preserved, got %q", result[0].Expected)
	}
}

func TestNormalizeFindings_PreservesService(t *testing.T) {
	findings := []Finding{
		normFinding("MyService", "image", "Redis:7", "redis:6"),
	}
	result := NormalizeFindings(findings, DefaultNormalizeOptions())
	if result[0].Service != "MyService" {
		t.Errorf("service name should not be modified, got %q", result[0].Service)
	}
}

func TestNormalizeEnvMap_TrimAndSort(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "  zval  ",
		"A_KEY": "\taval\t",
	}
	out, keys := NormalizeEnvMap(env, true)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "A_KEY" || keys[1] != "Z_KEY" {
		t.Errorf("keys not sorted: %v", keys)
	}
	if out["A_KEY"] != "aval" {
		t.Errorf("expected trimmed 'aval', got %q", out["A_KEY"])
	}
	if out["Z_KEY"] != "zval" {
		t.Errorf("expected trimmed 'zval', got %q", out["Z_KEY"])
	}
}

func TestNormalizeEnvMap_NoTrim(t *testing.T) {
	env := map[string]string{"KEY": "  value  "}
	out, _ := NormalizeEnvMap(env, false)
	if out["KEY"] != "  value  " {
		t.Errorf("expected untrimmed value, got %q", out["KEY"])
	}
}
