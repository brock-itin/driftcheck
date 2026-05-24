package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/driftcheck/internal/drift"
)

// WriteDiff writes a unified-diff-style output comparing expected vs actual
// values for each finding in the report.
func WriteDiff(w io.Writer, r drift.Report) error {
	if len(r.Findings) == 0 {
		_, err := fmt.Fprintln(w, "No drift detected — all services match their definitions.")
		return err
	}

	services := groupByService(r.Findings)
	for _, svc := range sortedKeys(services) {
		findings := services[svc]
		fmt.Fprintf(w, "--- expected/%s\n", svc)
		fmt.Fprintf(w, "+++ running/%s\n", svc)
		for _, f := range findings {
			writeFindingDiff(w, f)
		}
		fmt.Fprintln(w)
	}
	return nil
}

func writeFindingDiff(w io.Writer, f drift.Finding) {
	label := fmt.Sprintf("@@ %s @@", f.Field)
	fmt.Fprintln(w, label)
	if f.Expected != "" {
		fmt.Fprintf(w, "- %s\n", f.Expected)
	} else {
		fmt.Fprintln(w, "- (not set)")
	}
	if f.Actual != "" {
		fmt.Fprintf(w, "+ %s\n", f.Actual)
	} else {
		fmt.Fprintln(w, "+ (not set)")
	}
}

func groupByService(findings []drift.Finding) map[string][]drift.Finding {
	m := make(map[string][]drift.Finding)
	for _, f := range findings {
		m[f.Service] = append(m[f.Service], f)
	}
	return m
}

func sortedKeys(m map[string][]drift.Finding) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// simple insertion sort for deterministic output
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && strings.ToLower(keys[j]) < strings.ToLower(keys[j-1]); j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}
