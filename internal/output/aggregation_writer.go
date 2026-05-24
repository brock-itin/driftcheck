package output

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/yourusername/driftcheck/internal/drift"
)

// WriteAggregation writes a human-readable summary of aggregated drift windows.
func WriteAggregation(w io.Writer, result drift.AggregationResult, windowSize time.Duration) {
	if len(result.Windows) == 0 {
		fmt.Fprintln(w, "No aggregation data available.")
		return
	}

	fmt.Fprintf(w, "Drift Aggregation  (window: %s)\n", windowSize)
	fmt.Fprintln(w, strings.Repeat("-", 52))

	for i, win := range result.Windows {
		bar := buildBar(len(win.Findings), result.TotalFindings, 20)
		fmt.Fprintf(w, "[%2d] %s  %s  (%d)\n",
			i+1,
			win.Start.Format("2006-01-02 15:04"),
			bar,
			len(win.Findings),
		)
	}

	fmt.Fprintln(w, strings.Repeat("-", 52))
	fmt.Fprintf(w, "Total findings : %d\n", result.TotalFindings)

	if result.PeakWindow != nil {
		fmt.Fprintf(w, "Peak window    : %s  (%d findings)\n",
			result.PeakWindow.Start.Format("2006-01-02 15:04"),
			len(result.PeakWindow.Findings),
		)
	}
}

// AggregationSummaryLine returns a one-line summary string.
func AggregationSummaryLine(result drift.AggregationResult) string {
	if result.TotalFindings == 0 {
		return "aggregation: no findings recorded"
	}
	return fmt.Sprintf("aggregation: %d findings across %d window(s)",
		result.TotalFindings, len(result.Windows))
}

// buildBar constructs a simple ASCII bar proportional to count/total.
func buildBar(count, total, width int) string {
	if total == 0 || width == 0 {
		return strings.Repeat(".", width)
	}
	filled := count * width / total
	if filled == 0 && count > 0 {
		filled = 1
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}
