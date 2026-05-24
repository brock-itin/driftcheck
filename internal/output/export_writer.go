package output

import (
	"fmt"
	"io"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
)

// WriteExportSummary prints a human-readable summary of what was exported.
func WriteExportSummary(w io.Writer, records []drift.ExportRecord, format drift.ExportFormat, dest string) {
	if len(records) == 0 {
		fmt.Fprintln(w, "export: no findings to export")
		return
	}
	fmt.Fprintf(w, "export: wrote %d finding(s) as %s", len(records), format)
	if dest != "" {
		fmt.Fprintf(w, " → %s", dest)
	}
	fmt.Fprintln(w)
}

// ExportAndWrite is a convenience wrapper: converts the report, writes the
// export payload to payload, and writes a summary line to summary.
func ExportAndWrite(
	payload io.Writer,
	summary io.Writer,
	r drift.Report,
	format drift.ExportFormat,
	dest string,
	now time.Time,
) error {
	records := drift.ExportReport(r, now)
	if err := drift.WriteExport(payload, records, format); err != nil {
		return fmt.Errorf("export write: %w", err)
	}
	WriteExportSummary(summary, records, format, dest)
	return nil
}
