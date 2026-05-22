package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/driftcheck/internal/drift"
)

// Format controls the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes drift reports to an io.Writer in a specified format.
type Formatter struct {
	Writer io.Writer
	Format Format
}

// NewFormatter creates a Formatter with the given writer and format.
func NewFormatter(w io.Writer, f Format) *Formatter {
	return &Formatter{Writer: w, Format: f}
}

// Write renders the report to the formatter's writer.
func (f *Formatter) Write(r *drift.Report) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(r)
	default:
		return f.writeText(r)
	}
}

func (f *Formatter) writeJSON(r *drift.Report) error {
	enc := json.NewEncoder(f.Writer)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func (f *Formatter) writeText(r *drift.Report) error {
	summary := r.Summary()
	if summary.Total == 0 {
		fmt.Fprintln(f.Writer, "✓ No drift detected.")
		return nil
	}

	fmt.Fprintf(f.Writer, "Drift detected: %d finding(s)\n\n", summary.Total)

	tw := tabwriter.NewWriter(f.Writer, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tFIELD\tEXPECTED\tACTUAL")
	fmt.Fprintln(tw, "-------\t-----\t--------\t------")
	for _, finding := range r.Findings {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			finding.Service,
			finding.Field,
			finding.Expected,
			finding.Actual,
		)
	}
	return tw.Flush()
}
