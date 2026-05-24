package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/driftcheck/internal/drift"
)

// WriteThresholdResult writes a human-readable threshold evaluation result to w.
func WriteThresholdResult(w io.Writer, res drift.ThresholdResult) {
	if !res.Breached {
		fmt.Fprintln(w, "✓ All drift thresholds within limits.")
		return
	}

	action := strings.ToUpper(string(res.Action))
	fmt.Fprintf(w, "⚠ Threshold breach [%s]:\n", action)
	for _, msg := range res.Messages {
		fmt.Fprintf(w, "  - %s\n", msg)
	}
}

// ThresholdStatusLine returns a single-line summary suitable for CI output.
func ThresholdStatusLine(res drift.ThresholdResult) string {
	if !res.Breached {
		return "thresholds: OK"
	}
	return fmt.Sprintf("thresholds: BREACHED (%d violation(s)) [%s]",
		len(res.Messages),
		strings.ToUpper(string(res.Action)),
	)
}
