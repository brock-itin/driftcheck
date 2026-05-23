package snapshot_test

import (
	"errors"
	"testing"

	"github.com/yourorg/driftcheck/internal/docker"
	"github.com/yourorg/driftcheck/internal/snapshot"
)

// stubLister is a test double for docker.ContainerLister.
type stubLister struct {
	containers []docker.ContainerInfo
	err        error
}

func (s *stubLister) ListContainers() ([]docker.ContainerInfo, error) {
	return s.containers, s.err
}

// makeBuilder is a helper that constructs a Builder from a slice of ContainerInfo.
func makeBuilder(containers []docker.ContainerInfo) *snapshot.Builder {
	return snapshot.NewBuilder(&stubLister{containers: containers})
}

func TestBuild_PopulatesSnapshot(t *testing.T) {
	b := makeBuilder([]docker.ContainerInfo{
		{
			ID:      "c1",
			Name:    "web",
			Image:   "nginx:1.25",
			Env:     []string{"PORT=80", "DEBUG=false"},
			Labels:  map[string]string{"app": "web"},
			Running: true,
		},
	})
	snap, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	state, ok := snap.Containers["web"]
	if !ok {
		t.Fatal("expected 'web' in snapshot")
	}
	if state.Image != "nginx:1.25" {
		t.Errorf("expected nginx:1.25, got %s", state.Image)
	}
	if state.Env["PORT"] != "80" {
		t.Errorf("expected PORT=80, got %s", state.Env["PORT"])
	}
}

func TestBuild_EmptyContainers(t *testing.T) {
	b := makeBuilder([]docker.ContainerInfo{})
	snap, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Containers) != 0 {
		t.Errorf("expected empty snapshot, got %d entries", len(snap.Containers))
	}
}

func TestBuild_ListerError(t *testing.T) {
	lister := &stubLister{err: errors.New("docker unavailable")}
	b := snapshot.NewBuilder(lister)
	_, err := b.Build()
	if err == nil {
		t.Fatal("expected error from lister")
	}
}

func TestBuild_EnvWithoutValue(t *testing.T) {
	b := makeBuilder([]docker.ContainerInfo{
		{ID: "c2", Name: "worker", Image: "alpine", Env: []string{"FLAG"}, Running: true},
	})
	snap, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, ok := snap.Containers["worker"].Env["FLAG"]
	if !ok {
		t.Fatal("expected FLAG key in env")
	}
	if val != "" {
		t.Errorf("expected empty value for FLAG, got %q", val)
	}
}
