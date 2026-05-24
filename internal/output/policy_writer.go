package output

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/yourorg/driftcheck/internal/drift"
)

// WritePolicyEvaluation renders a table showing each finding and its resolved
// policy action to w.
func WritePolicyEvaluation(w io.Writer, findings []drift.Finding, policy *drift.Policy) error {
	if len(findings) == 0 {
		_, err := fmt.Fprintln(w, "No findings — policy evaluation skipped.")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tTYPE\tDRIFT\tACTION")
	fmt.Fprintln(tw, "-------\t----\t-----\t------")

	for _, f := range findings {
		action := policy.Resolve(f)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			f.Service,
			f.Type,
			truncate(f.Detail, 48),
			string(action),
		)
	}
	return tw.Flush()
}

// PolicyActionCounts returns a map of action -> count across all findings.
func PolicyActionCounts(findings []drift.Finding, policy *drift.Policy) map[drift.PolicyAction]int {
	counts := map[drift.PolicyAction]int{}
	for _, f := range findings {
		a := policy.Resolve(f)
		counts[a]++
	}
	return counts
}
