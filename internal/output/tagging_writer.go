package output

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/driftcheck/internal/drift"
)

// WriteTagged writes a human-readable table of tagged findings to w.
func WriteTagged(w io.Writer, tags map[string]drift.TagSet) {
	if len(tags) == 0 {
		fmt.Fprintln(w, "No tagged findings.")
		return
	}

	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "%-50s  %s\n", "FINDING", "TAGS")
	fmt.Fprintf(w, "%s  %s\n", strings.Repeat("-", 50), strings.Repeat("-", 30))
	for _, k := range keys {
		ts := tags[k]
		fmt.Fprintf(w, "%-50s  %s\n", truncateTag(k, 50), strings.Join(ts, ", "))
	}
	fmt.Fprintf(w, "\nTotal: %d finding(s) tagged\n", len(tags))
}

// TaggingStatusLine returns a one-line status suitable for CLI output.
func TaggingStatusLine(tags map[string]drift.TagSet) string {
	if len(tags) == 0 {
		return "tagging: no rules matched"
	}
	allTags := make(map[string]struct{})
	for _, ts := range tags {
		for _, t := range ts {
			allTags[t] = struct{}{}
		}
	}
	return fmt.Sprintf("tagging: %d finding(s) tagged with %d unique tag(s)", len(tags), len(allTags))
}

func truncateTag(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
