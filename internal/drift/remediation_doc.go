// Package drift provides drift detection, analysis, and remediation utilities
// for comparing running containers against their source compose definitions.
//
// # Remediation
//
// The Remediator type generates actionable remediation suggestions for each
// detected drift finding. It maps drift types (image, env, label) to concrete
// docker compose commands that restore the expected container state.
//
// Usage:
//
//	remediator := drift.NewRemediator("docker-compose.yml")
//	actions := remediator.Suggest(report.Findings)
//	output.WriteRemediation(os.Stdout, actions)
//
// Each RemediationAction includes:
//   - Service: the affected service name
//   - DriftType: the category of drift (image, env, label)
//   - Description: a human-readable explanation of the fix
//   - Command: the exact shell command to remediate the drift
package drift
