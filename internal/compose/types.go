// Package compose provides types and utilities for parsing Docker Compose files.
package compose

// Service represents a single service entry from a Docker Compose file.
// It captures the fields that driftcheck uses when comparing against running
// containers.
type Service struct {
	// Name is the service key as it appears in the compose file.
	Name string

	// Image is the fully-qualified image reference, e.g. "nginx:1.25".
	Image string

	// Environment holds key=value pairs declared under the service's
	// `environment` block. Both map and list forms are normalised into this
	// map by the parser.
	Environment map[string]string

	// Labels holds Docker labels declared under the service's `labels` block.
	// Both map and list forms are normalised into this map by the parser.
	Labels map[string]string

	// Ports contains the raw port-mapping strings, e.g. "8080:80".
	Ports []string

	// Volumes contains the raw volume-mount strings, e.g. "./data:/data".
	Volumes []string
}

// Compose is the top-level structure parsed from a docker-compose.yml file.
type Compose struct {
	// Version is the compose file format version string.
	Version string

	// Services maps service names to their definitions.
	Services map[string]Service
}

// ServiceNames returns a sorted list of service names defined in the compose
// file. The order is deterministic for consistent output.
func (c *Compose) ServiceNames() []string {
	if c == nil || len(c.Services) == 0 {
		return nil
	}
	names := make([]string, 0, len(c.Services))
	for name := range c.Services {
		names = append(names, name)
	}
	sortStrings(names)
	return names
}

// sortStrings sorts a string slice in-place (avoids importing sort in callers).
func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
