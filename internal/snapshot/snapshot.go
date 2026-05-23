// Package snapshot provides functionality to capture and persist
// the current state of running containers for later drift comparison.
package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// ContainerState represents the captured state of a single container.
type ContainerState struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	Env     map[string]string `json:"env"`
	Labels  map[string]string `json:"labels"`
	Running bool              `json:"running"`
}

// Snapshot holds a point-in-time capture of container states.
type Snapshot struct {
	CapturedAt time.Time                 `json:"captured_at"`
	Containers map[string]ContainerState `json:"containers"`
}

// New creates an empty Snapshot with the current timestamp.
func New() *Snapshot {
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Containers: make(map[string]ContainerState),
	}
}

// Add inserts or replaces a container state entry keyed by service name.
func (s *Snapshot) Add(serviceName string, state ContainerState) {
	s.Containers[serviceName] = state
}

// Save writes the snapshot as JSON to the given file path.
func (s *Snapshot) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, err
	}
	return &snap, nil
}
