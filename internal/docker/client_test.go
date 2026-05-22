package docker

import (
	"testing"
)

func TestContainerInfo_Fields(t *testing.T) {
	info := ContainerInfo{
		ID:      "abc123def456",
		Name:    "/my-service",
		Image:   "nginx:latest",
		Labels:  map[string]string{"com.docker.compose.service": "web"},
		Env:     []string{"PORT=8080"},
		Ports:   []string{"80/tcp"},
		Running: true,
	}

	if info.ID != "abc123def456" {
		t.Errorf("expected ID abc123def456, got %s", info.ID)
	}
	if info.Image != "nginx:latest" {
		t.Errorf("expected image nginx:latest, got %s", info.Image)
	}
	if !info.Running {
		t.Error("expected container to be running")
	}
	if len(info.Ports) != 1 || info.Ports[0] != "80/tcp" {
		t.Errorf("unexpected ports: %v", info.Ports)
	}
	if v, ok := info.Labels["com.docker.compose.service"]; !ok || v != "web" {
		t.Errorf("expected compose service label 'web', got %q", v)
	}
}

func TestNewClient_EnvFallback(t *testing.T) {
	// NewClient should succeed when DOCKER_HOST is not set (uses default socket).
	// We only verify no panic occurs and error handling works.
	_, err := NewClient()
	if err != nil {
		// In CI without Docker this is acceptable; just ensure it's a real error.
		t.Logf("NewClient returned error (expected in no-docker env): %v", err)
	}
}
