package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/driftcheck/internal/snapshot"
)

func makeState(image string) snapshot.ContainerState {
	return snapshot.ContainerState{
		ID:      "abc123",
		Name:    "web",
		Image:   image,
		Env:     map[string]string{"PORT": "8080"},
		Labels:  map[string]string{"app": "web"},
		Running: true,
	}
}

func TestNew_InitializesFields(t *testing.T) {
	s := snapshot.New()
	if s.Containers == nil {
		t.Fatal("expected Containers map to be initialized")
	}
	if s.CapturedAt.IsZero() {
		t.Fatal("expected CapturedAt to be set")
	}
}

func TestAdd_StoresState(t *testing.T) {
	s := snapshot.New()
	state := makeState("nginx:latest")
	s.Add("web", state)
	got, ok := s.Containers["web"]
	if !ok {
		t.Fatal("expected 'web' entry in snapshot")
	}
	if got.Image != "nginx:latest" {
		t.Errorf("expected image nginx:latest, got %s", got.Image)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	s := snapshot.New()
	s.CapturedAt = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	s.Add("web", makeState("nginx:1.25"))

	if err := s.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.CapturedAt != s.CapturedAt {
		t.Errorf("CapturedAt mismatch: got %v, want %v", loaded.CapturedAt, s.CapturedAt)
	}
	if loaded.Containers["web"].Image != "nginx:1.25" {
		t.Errorf("unexpected image: %s", loaded.Containers["web"].Image)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	s := snapshot.New()
	err := s.Save("/nonexistent/dir/snap.json")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestLoad_CorruptJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("{invalid json"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for corrupt JSON")
	}
}
