package drift

import (
	"strings"
	"testing"
)

func remFinding(service, driftType string) Finding {
	return Finding{Service: service, Type: driftType, Field: "test", Expected: "a", Actual: "b"}
}

func TestNewRemediator_SetsComposeFile(t *testing.T) {
	r := NewRemediator("docker-compose.yml")
	if r.ComposeFile != "docker-compose.yml" {
		t.Errorf("expected docker-compose.yml, got %s", r.ComposeFile)
	}
}

func TestSuggest_ImageDrift(t *testing.T) {
	r := NewRemediator("docker-compose.yml")
	actions := r.Suggest([]Finding{remFinding("web", "image")})
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if !strings.Contains(actions[0].Command, "up -d --no-deps web") {
		t.Errorf("unexpected command: %s", actions[0].Command)
	}
	if actions[0].Service != "web" {
		t.Errorf("expected service web, got %s", actions[0].Service)
	}
}

func TestSuggest_EnvDrift(t *testing.T) {
	r := NewRemediator("compose.yml")
	actions := r.Suggest([]Finding{remFinding("api", "env")})
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if !strings.Contains(actions[0].Command, "--force-recreate") {
		t.Errorf("expected force-recreate in command: %s", actions[0].Command)
	}
}

func TestSuggest_LabelDrift(t *testing.T) {
	r := NewRemediator("compose.yml")
	actions := r.Suggest([]Finding{remFinding("worker", "label")})
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if !strings.Contains(actions[0].Command, "worker") {
		t.Errorf("expected worker in command: %s", actions[0].Command)
	}
}

func TestSuggest_UnknownType_Skipped(t *testing.T) {
	r := NewRemediator("compose.yml")
	actions := r.Suggest([]Finding{remFinding("svc", "unknown")})
	if len(actions) != 0 {
		t.Errorf("expected 0 actions for unknown type, got %d", len(actions))
	}
}

func TestSuggest_MultipleFindings(t *testing.T) {
	r := NewRemediator("docker-compose.yml")
	findings := []Finding{
		remFinding("web", "image"),
		remFinding("api", "env"),
		remFinding("db", "unknown"),
	}
	actions := r.Suggest(findings)
	if len(actions) != 2 {
		t.Errorf("expected 2 actions, got %d", len(actions))
	}
}

func TestSuggest_EmptyFindings(t *testing.T) {
	r := NewRemediator("compose.yml")
	actions := r.Suggest(nil)
	if len(actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(actions))
	}
}
