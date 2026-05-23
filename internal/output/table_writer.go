package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/yourusername/driftcheck/internal/drift"
)

const (
	minWidth = 0
	tabWidth = 0
	padding  = 2
	padChar  = ' '
)

// WriteTable writes drift findings as an aligned table to w.
// Each row contains: SERVICE, TYPE, FIELD, EXPECTED, ACTUAL, SEVERITY.
func WriteTable(w io.Writer, r drift.Report) error {
	tw := tabwriter.NewWriter(w, minWidth, tabWidth, padding, padChar, 0)

	if len(r.Findings) == 0 {
		_, err := fmt.Fprintln(w, "No drift detected.")
		return err
	}

	header := strings.Join([]string{"SERVICE", "TYPE", "FIELD", "EXPECTED", "ACTUAL", "SEVERITY"}, "\t")
	if _, err := fmt.Fprintln(tw, header); err != nil {
		return err
	}

	sep := strings.Repeat("-", 72)
	if _, err := fmt.Fprintln(tw, sep); err != nil {
		return err
	}

	for _, f := range r.Findings {
		row := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s",
			f.Service,
			f.Type,
			f.Field,
			truncate(f.Expected, 30),
			truncate(f.Actual, 30),
			f.Severity.String(),
		)
		if _, err := fmt.Fprintln(tw, row); err != nil {
			return err
		}
	}

	return tw.Flush()
}

// truncate shortens s to max runes, appending "…" if trimmed.
func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "…"
}
