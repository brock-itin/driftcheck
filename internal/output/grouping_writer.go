package output

import (
	"fmt"
	"io"

	"github.com/yourorg/driftcheck/internal/drift"
)

// WriteGroupedByType writes findings to w organised by drift type.
func WriteGroupedByType(w io.Writer, r drift.Report) {
	groups := drift.GroupByType(r)
	if len(groups) == 0 {
		fmt.Fprintln(w, "No drift findings.")
		return
	}
	for _, typ := range drift.SortedGroupKeys(groups) {
		findings := groups[typ]
		fmt.Fprintf(w, "\n[type: %s] (%d finding(s))\n", typ, len(findings))
		for _, f := range findings {
			fmt.Fprintf(w, "  service=%-20s field=%-20s want=%s got=%s\n",
				f.Service, f.Field, f.Want, f.Got)
		}
	}
}

// WriteGroupedByService writes findings to w organised by service name.
func WriteGroupedByService(w io.Writer, r drift.Report) {
	groups := drift.GroupByService(r)
	if len(groups) == 0 {
		fmt.Fprintln(w, "No drift findings.")
		return
	}
	for _, svc := range drift.SortedGroupKeys(groups) {
		findings := groups[svc]
		fmt.Fprintf(w, "\n[service: %s] (%d finding(s))\n", svc, len(findings))
		for _, f := range findings {
			fmt.Fprintf(w, "  type=%-12s field=%-20s want=%s got=%s\n",
				f.Type, f.Field, f.Want, f.Got)
		}
	}
}

// GroupingSummaryLine returns a one-line summary of grouped findings.
func GroupingSummaryLine(r drift.Report) string {
	byType := drift.GroupByType(r)
	if len(byType) == 0 {
		return "grouped summary: no drift detected"
	}
	return fmt.Sprintf("grouped summary: %d type(s) across %d finding(s)",
		len(byType), len(r.Findings))
}
