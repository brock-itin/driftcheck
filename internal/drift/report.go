package drift

import (
	"fmt"
	"io"
	"strings"
)

// ReportFormat controls the output format of a drift report.
type ReportFormat string

const (
	FormatText ReportFormat = "text"
	FormatJSON ReportFormat = "json"
)

// Report writes a human-readable summary of drift findings to w.
func Report(w io.Writer, findings []Finding) {
	if len(findings) == 0 {
		fmt.Fprintln(w, "✓ No drift detected. All services match their compose definitions.")
		return
	}

	fmt.Fprintf(w, "⚠ Drift detected: %d finding(s)\n", len(findings))
	fmt.Fprintln(w, strings.Repeat("-", 60))

	for i, f := range findings {
		fmt.Fprintf(w, "[%d] Service : %s\n", i+1, f.Service)
		fmt.Fprintf(w, "    Type    : %s\n", f.Type)
		fmt.Fprintf(w, "    Expected: %s\n", f.Expected)
		fmt.Fprintf(w, "    Actual  : %s\n", f.Actual)
		fmt.Fprintf(w, "    Detail  : %s\n", f.Message)
		if i < len(findings)-1 {
			fmt.Fprintln(w, strings.Repeat("-", 60))
		}
	}
}

// Summary returns a one-line summary string suitable for logging or CI output.
func Summary(findings []Finding) string {
	if len(findings) == 0 {
		return "no drift detected"
	}
	services := make(map[string]struct{})
	for _, f := range findings {
		services[f.Service] = struct{}{}
	}
	names := make([]string, 0, len(services))
	for s := range services {
		names = append(names, s)
	}
	return fmt.Sprintf("%d finding(s) across %d service(s): %s",
		len(findings), len(services), strings.Join(names, ", "))
}
