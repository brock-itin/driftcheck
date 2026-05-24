package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/you/driftcheck/internal/drift"
)

// WriteTrend writes a human-readable trend summary to w.
func WriteTrend(w io.Writer, t drift.Trend) {
	if len(t.Points) == 0 {
		fmt.Fprintln(w, "No trend data available.")
		return
	}
	fmt.Fprintf(w, "Drift Trend (%d snapshots)\n", len(t.Points))
	fmt.Fprintln(w, strings.Repeat("-", 50))
	for _, p := range t.Points {
		ts := p.Timestamp.Format("2006-01-02 15:04:05")
		fmt.Fprintf(w, "  %s  total=%-4d", ts, p.Total)
		for _, k := range sortedMapKeys(p.ByType) {
			fmt.Fprintf(w, "  %s=%d", k, p.ByType[k])
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w, strings.Repeat("-", 50))
	delta := t.Delta()
	direction := "unchanged"
	switch {
	case delta > 0:
		direction = fmt.Sprintf("▲ +%d", delta)
	case delta < 0:
		direction = fmt.Sprintf("▼ %d", delta)
	}
	fmt.Fprintf(w, "Latest: %d findings  (%s since previous)\n", t.Latest().Total, direction)
}

// TrendDeltaLine returns a compact one-line summary of the trend delta.
func TrendDeltaLine(t drift.Trend) string {
	if len(t.Points) == 0 {
		return "no data"
	}
	delta := t.Delta()
	latest := t.Latest().Total
	switch {
	case delta > 0:
		return fmt.Sprintf("%d findings (+%d)", latest, delta)
	case delta < 0:
		return fmt.Sprintf("%d findings (%d)", latest, delta)
	default:
		return fmt.Sprintf("%d findings (no change)", latest)
	}
}
