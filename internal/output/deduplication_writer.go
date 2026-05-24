package output

import (
	"fmt"
	"io"

	"github.com/your-org/driftcheck/internal/drift"
)

// WriteDeduplicationSummary writes a human-readable summary of deduplication
// results to w, showing how many duplicates were removed and what remains.
func WriteDeduplicationSummary(w io.Writer, original, deduped drift.Report) {
	removed := drift.DuplicateCount(original, deduped)
	total := len(original.Findings)
	kept := len(deduped.Findings)

	fmt.Fprintf(w, "Deduplication Summary\n")
	fmt.Fprintf(w, "  Total findings (before): %d\n", total)
	fmt.Fprintf(w, "  Duplicates removed:      %d\n", removed)
	fmt.Fprintf(w, "  Unique findings (after): %d\n", kept)

	if removed == 0 {
		fmt.Fprintf(w, "  No duplicates detected.\n")
		return
	}

	fmt.Fprintf(w, "\n  Retained findings:\n")
	for _, f := range deduped.Findings {
		fmt.Fprintf(w, "    [%s] %s / %s: %s\n",
			f.Severity, f.Service, f.Type, f.Field)
	}
}

// DeduplicationStatusLine returns a compact single-line summary suitable for
// use in CI output or log lines.
func DeduplicationStatusLine(original, deduped drift.Report) string {
	removed := drift.DuplicateCount(original, deduped)
	if removed == 0 {
		return fmt.Sprintf("dedup: no duplicates in %d findings", len(original.Findings))
	}
	return fmt.Sprintf("dedup: removed %d duplicate(s), %d unique finding(s) remain",
		removed, len(deduped.Findings))
}
