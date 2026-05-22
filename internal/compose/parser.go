package compose

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceConfig represents a single service definition from a Docker Compose file.
type ServiceConfig struct {
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"environment"`
	Ports       []string          `yaml:"ports"`
	Volumes     []string          `yaml:"volumes"`
	Command     string            `yaml:"command"`
}

// ComposeFile represents the top-level structure of a docker-compose.yml.
type ComposeFile struct {
	Version  string                    `yaml:"version"`
	Services map[string]ServiceConfig  `yaml:"services"`
}

// ParseFile reads and parses a Docker Compose YAML file from the given path.
func ParseFile(path string) (*ComposeFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading compose file %q: %w", path, err)
	}

	var cf ComposeFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("parsing compose file %q: %w", path, err)
	}

	if len(cf.Services) == 0 {
		return nil, fmt.Errorf("compose file %q defines no services", path)
	}

	return &cf, nil
}

// ServiceNames returns the list of service names defined in the compose file.
func (cf *ComposeFile) ServiceNames() []string {
	names := make([]string, 0, len(cf.Services))
	for name := range cf.Services {
		names = append(names, name)
	}
	return names
}
