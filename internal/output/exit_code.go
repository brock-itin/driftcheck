package output

import "github.com/user/driftcheck/internal/drift"

// ExitCode constants for CLI exit status.
const (
	ExitOK      = 0
	ExitDrift   = 1
	ExitError   = 2
)

// ResolveExitCode returns the appropriate exit code based on the drift report.
// ExitOK is returned when no findings are present.
// ExitDrift is returned when one or more drift findings exist.
func ResolveExitCode(r drift.Report) int {
	if len(r.Findings) == 0 {
		return ExitOK
	}
	return ExitDrift
}
