package drift

import (
	"testing"

	"github.com/user/driftcheck/internal/docker"
)

func TestDetect_NoFindings(t *testing.T) {
	expected := map[string]ComposeService{
		"web": {Image: "nginx:1.25", Environment: map[string]string{"PORT": "8080"}},
	}
	running := map[string]docker.ContainerInfo{
		"web": {Image: "nginx:1.25", Env: map[string]string{"PORT": "8080"}},
	}

	findings := Detect(expected, running)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %+v", len(findings), findings)
	}
}

func TestDetect_ImageDrift(t *testing.T) {
	expected := map[string]ComposeService{
		"api": {Image: "myapp:v2"},
	}
	running := map[string]docker.ContainerInfo{
		"api": {Image: "myapp:v1", Env: map[string]string{}},
	}

	findings := Detect(expected, running)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Type != DriftImage {
		t.Errorf("expected DriftImage, got %s", findings[0].Type)
	}
	if findings[0].Expected != "myapp:v2" {
		t.Errorf("unexpected expected value: %s", findings[0].Expected)
	}
}

func TestDetect_EnvDrift(t *testing.T) {
	expected := map[string]ComposeService{
		"worker": {
			Image:       "worker:latest",
			Environment: map[string]string{"LOG_LEVEL": "info"},
		},
	}
	running := map[string]docker.ContainerInfo{
		"worker": {Image: "worker:latest", Env: map[string]string{"LOG_LEVEL": "debug"}},
	}

	findings := Detect(expected, running)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Type != DriftEnv {
		t.Errorf("expected DriftEnv, got %s", findings[0].Type)
	}
}

func TestDetect_MissingContainer(t *testing.T) {
	expected := map[string]ComposeService{
		"db": {Image: "postgres:15"},
	}
	running := map[string]docker.ContainerInfo{}

	findings := Detect(expected, running)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Service != "db" {
		t.Errorf("expected service 'db', got %s", findings[0].Service)
	}
}

func TestDetect_MultipleDrifts(t *testing.T) {
	expected := map[string]ComposeService{
		"web": {
			Image:       "nginx:1.25",
			Environment: map[string]string{"ENV": "prod", "DEBUG": "false"},
		},
	}
	running := map[string]docker.ContainerInfo{
		"web": {Image: "nginx:1.24", Env: map[string]string{"ENV": "staging", "DEBUG": "false"}},
	}

	findings := Detect(expected, running)
	// expect image drift + ENV drift
	if len(findings) != 2 {
		t.Errorf("expected 2 findings, got %d: %+v", len(findings), findings)
	}
}
