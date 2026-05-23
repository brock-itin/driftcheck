package snapshot

import (
	"strings"

	"github.com/yourorg/driftcheck/internal/docker"
)

// Builder constructs a Snapshot from live Docker container data.
type Builder struct {
	client docker.ContainerLister
}

// NewBuilder returns a Builder using the provided container lister.
func NewBuilder(client docker.ContainerLister) *Builder {
	return &Builder{client: client}
}

// Build queries running containers and returns a populated Snapshot.
func (b *Builder) Build() (*Snapshot, error) {
	containers, err := b.client.ListContainers()
	if err != nil {
		return nil, err
	}

	snap := New()
	for _, c := range containers {
		state := ContainerState{
			ID:      c.ID,
			Name:    c.Name,
			Image:   c.Image,
			Env:     parseEnvSlice(c.Env),
			Labels:  c.Labels,
			Running: c.Running,
		}
		snap.Add(c.Name, state)
	}
	return snap, nil
}

// parseEnvSlice converts ["KEY=VALUE", ...] into a map.
func parseEnvSlice(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		} else if len(parts) == 1 && parts[0] != "" {
			m[parts[0]] = ""
		}
	}
	return m
}
