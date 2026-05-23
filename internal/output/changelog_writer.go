package output

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
)

// WriteChangelog renders a Changelog to w in a human-readable table.
func WriteChangelog(w io.Writer, cl *drift.Changelog) error {
	if len(cl.Entries) == 0 {
		_, err := fmt.Fprintln(w, "No changelog entries recorded.")
		return err
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tCOMPOSE FILE\tTOTAL\tHIGH\tMEDIUM\tLOW")
	for _, e := range cl.Entries {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%d\t%d\t%d\n",
			e.Timestamp.Format(time.RFC3339),
			e.ComposFile,
			e.Total,
			e.BySeverity["high"],
			e.BySeverity["medium"],
			e.BySeverity["low"],
		)
	}
	return tw.Flush()
}

// ChangelogSummary returns a one-line summary of the most recent changelog entry.
func ChangelogSummary(cl *drift.Changelog) string {
	if len(cl.Entries) == 0 {
		return "no previous runs recorded"
	}
	last := cl.Entries[len(cl.Entries)-1]
	return fmt.Sprintf("last run: %s — %d finding(s) in %s",
		last.Timestamp.Format(time.RFC3339),
		last.Total,
		last.ComposFile,
	)
}
