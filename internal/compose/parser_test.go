package compose_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/driftcheck/internal/compose"
)

func writeTempCompose(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "docker-compose.yml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp compose file: %v", err)
	}
	return path
}

func TestParseFile_Valid(t *testing.T) {
	content := `
version: "3.9"
services:
  web:
    image: nginx:latest
    ports:
      - "80:80"
    environment:
      ENV: production
  db:
    image: postgres:14
    environment:
      POSTGRES_PASSWORD: secret
`
	path := writeTempCompose(t, content)

	cf, err := compose.ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cf.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(cf.Services))
	}

	web, ok := cf.Services["web"]
	if !ok {
		t.Fatal("expected service 'web' to exist")
	}
	if web.Image != "nginx:latest" {
		t.Errorf("expected image 'nginx:latest', got %q", web.Image)
	}
}

func TestParseFile_NoServices(t *testing.T) {
	content := `version: "3.9"
services: {}
`
	path := writeTempCompose(t, content)

	_, err := compose.ParseFile(path)
	if err == nil {
		t.Fatal("expected error for empty services, got nil")
	}
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := compose.ParseFile("/nonexistent/docker-compose.yml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestServiceNames(t *testing.T) {
	content := `
version: "3.9"
services:
  alpha:
    image: alpine
  beta:
    image: busybox
`
	path := writeTempCompose(t, content)
	cf, err := compose.ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	names := cf.ServiceNames()
	if len(names) != 2 {
		t.Errorf("expected 2 service names, got %d", len(names))
	}
}
