package drift

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"
)

// ExportFormat specifies the output format for drift exports.
type ExportFormat string

const (
	ExportJSON ExportFormat = "json"
	ExportCSV  ExportFormat = "csv"
)

// ExportRecord is a flat, serialisable representation of a single finding.
type ExportRecord struct {
	Timestamp string `json:"timestamp" csv:"timestamp"`
	Service   string `json:"service"   csv:"service"`
	Type      string `json:"type"      csv:"type"`
	Severity  string `json:"severity"  csv:"severity"`
	Expected  string `json:"expected"  csv:"expected"`
	Actual    string `json:"actual"    csv:"actual"`
	Message   string `json:"message"   csv:"message"`
}

// ExportReport converts a Report into a slice of ExportRecords stamped with now.
func ExportReport(r Report, now time.Time) []ExportRecord {
	ts := now.UTC().Format(time.RFC3339)
	records := make([]ExportRecord, 0, len(r.Findings))
	for _, f := range r.Findings {
		records = append(records, ExportRecord{
			Timestamp: ts,
			Service:   f.Service,
			Type:      string(f.Type),
			Severity:  f.Severity.String(),
			Expected:  f.Expected,
			Actual:    f.Actual,
			Message:   f.Message,
		})
	}
	sort.Slice(records, func(i, j int) bool {
		if records[i].Service != records[j].Service {
			return records[i].Service < records[j].Service
		}
		return records[i].Type < records[j].Type
	})
	return records
}

// WriteExport serialises records to w in the requested format.
func WriteExport(w io.Writer, records []ExportRecord, format ExportFormat) error {
	switch format {
	case ExportJSON:
		return writeExportJSON(w, records)
	case ExportCSV:
		return writeExportCSV(w, records)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}

func writeExportJSON(w io.Writer, records []ExportRecord) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func writeExportCSV(w io.Writer, records []ExportRecord) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"timestamp", "service", "type", "severity", "expected", "actual", "message"}); err != nil {
		return err
	}
	for _, r := range records {
		if err := cw.Write([]string{r.Timestamp, r.Service, r.Type, r.Severity, r.Expected, r.Actual, r.Message}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
