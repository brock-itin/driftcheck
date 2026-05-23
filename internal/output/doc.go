// Package output provides formatting, rendering, and exit-code resolution
// for driftcheck results.
//
// Formatters support multiple output modes:
//   - text: human-readable table-style output (default)
//   - json: machine-readable JSON suitable for CI pipelines
//
// Exit codes follow UNIX conventions:
//   - 0 (ExitOK):    no drift detected
//   - 1 (ExitDrift): one or more drift findings detected
//   - 2 (ExitError): a runtime or configuration error occurred
package output
