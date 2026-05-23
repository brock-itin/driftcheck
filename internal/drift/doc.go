// Package drift provides drift detection between running Docker containers
// and their source compose definitions.
//
// Core components:
//
//   - Detect: compares a compose file's service definitions against a live
//     container snapshot and returns a Report of Findings.
//
//   - Filter: narrows a Report by drift type, service name, or severity.
//
//   - Watcher: polls for drift at a configurable interval and invokes a
//     callback whenever new findings are detected.
//
//   - Baseline / Suppression / IgnoreList: mechanisms for acknowledging known
//     drift so that repeated alerts are suppressed.
//
//   - Changelog: append-only log of drift reports over time.
package drift
