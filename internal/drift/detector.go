package drift

import (
	"fmt"

	"github.com/user/driftcheck/internal/docker"
)

// DriftType categorizes the kind of drift detected.
type DriftType string

const (
	DriftImage DriftType = "image"
	DriftEnv   DriftType = "env"
	DriftPort  DriftType = "port"
)

// Finding represents a single drift finding for a service.
type Finding struct {
	Service  string
	Type     DriftType
	Expected string
	Actual   string
	Message  string
}

// ComposeService holds the expected state of a service from a compose file.
type ComposeService struct {
	Image       string
	Environment map[string]string
	Ports       []string
}

// Detect compares a map of expected compose services against running container
// info and returns all drift findings.
func Detect(expected map[string]ComposeService, running map[string]docker.ContainerInfo) []Finding {
	var findings []Finding

	for name, svc := range expected {
		container, ok := running[name]
		if !ok {
			findings = append(findings, Finding{
				Service:  name,
				Type:     DriftImage,
				Expected: svc.Image,
				Actual:   "",
				Message:  fmt.Sprintf("service %q not found in running containers", name),
			})
			continue
		}

		if svc.Image != "" && svc.Image != container.Image {
			findings = append(findings, Finding{
				Service:  name,
				Type:     DriftImage,
				Expected: svc.Image,
				Actual:   container.Image,
				Message:  fmt.Sprintf("image mismatch for service %q", name),
			})
		}

		for k, v := range svc.Environment {
			actual, exists := container.Env[k]
			if !exists || actual != v {
				findings = append(findings, Finding{
					Service:  name,
					Type:     DriftEnv,
					Expected: fmt.Sprintf("%s=%s", k, v),
					Actual:   fmt.Sprintf("%s=%s", k, actual),
					Message:  fmt.Sprintf("env var %q mismatch for service %q", k, name),
				})
			}
		}
	}

	return findings
}
