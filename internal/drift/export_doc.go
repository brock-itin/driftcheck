// Package drift — export.go
//
// # Drift Export
//
// The export subsystem converts a [Report] into flat, serialisable records
// suitable for archiving, feeding into external dashboards, or piping to other
// CLI tools.
//
// ## Formats
//
//   - ExportJSON — pretty-printed JSON array of [ExportRecord] objects.
//   - ExportCSV  — RFC 4180 CSV with a header row.
//
// ## Usage
//
//	records := drift.ExportReport(report, time.Now())
//	if err := drift.WriteExport(os.Stdout, records, drift.ExportJSON); err != nil {
//		log.Fatal(err)
//	}
//
// Records are deterministically sorted by (service, type) so that diff-friendly
// output is produced across successive runs against the same compose file.
package drift
