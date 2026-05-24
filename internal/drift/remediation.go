package drift

import (
	"fmt"
	"strings"
)

// RemediationAction describes a suggested fix for a drift finding.
type RemediationAction struct {
	Service     string
	DriftType   string
	Description string
	Command     string
}

// Remediator generates remediation suggestions for drift findings.
type Remediator struct {
	ComposeFile string
}

// NewRemediator creates a Remediator tied to the given compose file path.
func NewRemediator(composeFile string) *Remediator {
	return &Remediator{ComposeFile: composeFile}
}

// Suggest returns a list of remediation actions for the given findings.
func (r *Remediator) Suggest(findings []Finding) []RemediationAction {
	actions := make([]RemediationAction, 0, len(findings))
	for _, f := range findings {
		if a, ok := r.actionFor(f); ok {
			actions = append(actions, a)
		}
	}
	return actions
}

func (r *Remediator) actionFor(f Finding) (RemediationAction, bool) {
	switch strings.ToLower(f.Type) {
	case "image":
		return RemediationAction{
			Service:     f.Service,
			DriftType:   f.Type,
			Description: fmt.Sprintf("Re-deploy '%s' to restore expected image", f.Service),
			Command:     fmt.Sprintf("docker compose -f %s up -d --no-deps %s", r.ComposeFile, f.Service),
		}, true
	case "env":
		return RemediationAction{
			Service:     f.Service,
			DriftType:   f.Type,
			Description: fmt.Sprintf("Restart '%s' to apply expected environment variables", f.Service),
			Command:     fmt.Sprintf("docker compose -f %s up -d --force-recreate --no-deps %s", r.ComposeFile, f.Service),
		}, true
	case "label":
		return RemediationAction{
			Service:     f.Service,
			DriftType:   f.Type,
			Description: fmt.Sprintf("Recreate '%s' to restore expected labels", f.Service),
			Command:     fmt.Sprintf("docker compose -f %s up -d --force-recreate --no-deps %s", r.ComposeFile, f.Service),
		}, true
	}
	return RemediationAction{}, false
}
