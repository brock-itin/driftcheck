// Package output provides formatters for rendering drift reports
// to various output formats such as plain text tables and JSON.
//
// Usage:
//
//	f := output.NewFormatter(os.Stdout, output.FormatText)
//	f.Write(report)
//
// Supported formats:
//   - FormatText: human-readable tabular output
//   - FormatJSON: machine-readable JSON suitable for CI pipelines
package output
