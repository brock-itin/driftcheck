package output

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
)

// WriteEnriched renders enriched findings as a formatted table.
func WriteEnriched(w io.Writer, findings []drift.EnrichedFinding) error {
	if len(findings) == 0 {
		_, err := fmt.Fprintln(w, "no enriched findings")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tTYPE\tSEVERITY\tENVIRONMENT\tCLUSTER\tRUNBOOK\tENRICHED_AT")
	fmt.Fprintln(tw, "-------\t----\t--------\t-----------\t-------\t-------\t-----------")

	for _, ef := range findings {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			ef.Service,
			ef.Type,
			ef.Severity,
			envOrDash(ef.Environment),
			envOrDash(ef.Cluster),
			runbookOrDash(ef.RunbookURL),
			ef.EnrichedAt.Format(time.RFC3339),
		)
	}
	return tw.Flush()
}

// EnrichmentStatusLine returns a one-line summary suitable for CLI output.
func EnrichmentStatusLine(findings []drift.EnrichedFinding, src drift.EnrichmentSource) string {
	if len(findings) == 0 {
		return "enrichment complete: no findings"
	}
	env := src.Environment
	if env == "" {
		env = "(unset)"
	}
	return fmt.Sprintf("enrichment complete: %d finding(s) tagged [env=%s cluster=%s]",
		len(findings), env, envOrDash(src.Cluster))
}

func envOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func runbookOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
