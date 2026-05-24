package output

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/yourorg/driftcheck/internal/drift"
)

// WriteMetrics renders a Metrics summary to w in human-readable form.
func WriteMetrics(w io.Writer, m drift.Metrics) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Run at:\t%s\n", m.RunAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(tw, "Duration:\t%dms\n", m.DurationMs)
	fmt.Fprintf(tw, "Services scanned:\t%d\n", m.ServicesTotal)
	fmt.Fprintf(tw, "Services with drift:\t%d\n", m.ServicesDrifted)
	fmt.Fprintf(tw, "Total findings:\t%d\n", m.FindingsTotal)
	_ = tw.Flush()

	if len(m.BySeverity) > 0 {
		fmt.Fprintln(w, "\nFindings by severity:")
		keys := sortedMapKeys(m.BySeverity)
		for _, k := range keys {
			fmt.Fprintf(w, "  %-10s %d\n", k, m.BySeverity[k])
		}
	}

	if len(m.ByType) > 0 {
		fmt.Fprintln(w, "\nFindings by type:")
		keys := sortedMapKeys(m.ByType)
		for _, k := range keys {
			fmt.Fprintf(w, "  %-10s %d\n", k, m.ByType[k])
		}
	}
}

func sortedMapKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// MetricsSummaryLine returns a compact one-line summary string.
func MetricsSummaryLine(m drift.Metrics) string {
	return fmt.Sprintf(
		"scanned %d service(s) in %dms — %d finding(s) across %d service(s)",
		m.ServicesTotal, m.DurationMs, m.FindingsTotal, m.ServicesDrifted,
	)
}
