// Package drift provides drift detection, filtering, severity annotation,
// suppression, ignore lists, baselines, changelogs, and policy enforcement
// for driftcheck.
//
// # Policy
//
// A policy file (YAML) defines ordered rules that map drift types and optional
// service names to one of three actions:
//
//   - warn   – report the finding but exit 0 (default)
//   - fail   – report the finding and exit non-zero
//   - ignore – silently drop the finding
//
// Example policy.yaml:
//
//	rules:
//	  - type: image
//	    action: fail
//	  - type: env
//	    service: legacy-worker
//	    action: ignore
//	  - type: "*"
//	    action: warn
//
// Rules are evaluated top-to-bottom; the first match wins.
// If no rule matches, the default action is warn.
package drift
