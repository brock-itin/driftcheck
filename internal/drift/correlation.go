package drift

import (
	"sort"
	"time"
)

// CorrelationWindow groups findings that occurred within a shared time window
// and share a common service or type, suggesting a correlated root cause.
type CorrelationWindow struct {
	WindowStart time.Time
	WindowEnd   time.Time
	Findings    []Finding
	Services    []string
	Types       []string
}

// CorrelationOptions controls how findings are correlated.
type CorrelationOptions struct {
	// Window is the maximum time gap between findings to consider them correlated.
	Window time.Duration
	// MinSize is the minimum number of findings required to form a correlation group.
	MinSize int
}

// DefaultCorrelationOptions returns sensible defaults.
func DefaultCorrelationOptions() CorrelationOptions {
	return CorrelationOptions{
		Window:  5 * time.Minute,
		MinSize: 2,
	}
}

// CorrelateFindings groups findings into correlation windows based on timestamp
// proximity and shared service or type attributes.
func CorrelateFindings(findings []Finding, opts CorrelationOptions) []CorrelationWindow {
	if len(findings) == 0 {
		return nil
	}
	if opts.Window <= 0 || opts.MinSize < 1 {
		return nil
	}

	sorted := make([]Finding, len(findings))
	copy(sorted, findings)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].DetectedAt.Before(sorted[j].DetectedAt)
	})

	var windows []CorrelationWindow
	i := 0
	for i < len(sorted) {
		start := sorted[i].DetectedAt
		end := start.Add(opts.Window)
		var group []Finding
		for j := i; j < len(sorted) && !sorted[j].DetectedAt.After(end); j++ {
			group = append(group, sorted[j])
		}
		if len(group) >= opts.MinSize {
			windows = append(windows, buildCorrelationWindow(group, start, end))
		}
		i += len(group)
		if len(group) == 0 {
			i++
		}
	}
	return windows
}

func buildCorrelationWindow(group []Finding, start, end time.Time) CorrelationWindow {
	svcSet := map[string]struct{}{}
	typeSet := map[string]struct{}{}
	for _, f := range group {
		svcSet[f.Service] = struct{}{}
		typeSet[f.Type] = struct{}{}
	}
	return CorrelationWindow{
		WindowStart: start,
		WindowEnd:   end,
		Findings:    group,
		Services:    setToSortedSlice(svcSet),
		Types:       setToSortedSlice(typeSet),
	}
}

func setToSortedSlice(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
