package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/example/driftcheck/internal/drift"
)

// WriteSummaryTable writes a grouped summary of drift findings by service and type.
func WriteSummaryTable(w io.Writer, report drift.Report) {
	if len(report.Findings) == 0 {
		fmt.Fprintln(w, "No drift detected across all services.")
		return
	}

	type key struct {
		Service string
		Type    string
	}

	counts := make(map[key]int)
	services := make(map[string]struct{})

	for _, f := range report.Findings {
		k := key{Service: f.Service, Type: f.Type}
		counts[k]++
		services[f.Service] = struct{}{}
	}

	fmt.Fprintf(w, "Drift Summary (%d finding(s) across %d service(s))\n", len(report.Findings), len(services))
	fmt.Fprintln(w, strings.Repeat("-", 52))
	fmt.Fprintf(w, "%-24s %-16s %s\n", "SERVICE", "TYPE", "COUNT")
	fmt.Fprintln(w, strings.Repeat("-", 52))

	for k, count := range counts {
		fmt.Fprintf(w, "%-24s %-16s %d\n", truncate(k.Service, 22), truncate(k.Type, 14), count)
	}

	fmt.Fprintln(w, strings.Repeat("-", 52))
}

// ServiceDriftStatus returns a map of service name to whether drift was found.
func ServiceDriftStatus(report drift.Report) map[string]bool {
	status := make(map[string]bool)
	for _, f := range report.Findings {
		if !status[f.Service] {
			status[f.Service] = true
		}
	}
	return status
}
