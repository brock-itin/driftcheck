package drift

import (
	"sort"
	"time"
)

// AggregationWindow defines a time bucket for aggregating findings.
type AggregationWindow struct {
	Start    time.Time
	End      time.Time
	Findings []Finding
}

// AggregationResult holds bucketed findings and summary stats.
type AggregationResult struct {
	Windows      []AggregationWindow
	TotalFindings int
	PeakWindow   *AggregationWindow
}

// AggregateByWindow groups changelog entries into fixed-duration time windows
// and returns the findings that appeared in each bucket.
func AggregateByWindow(entries []ChangelogEntry, windowSize time.Duration) AggregationResult {
	if len(entries) == 0 || windowSize <= 0 {
		return AggregationResult{}
	}

	// Sort entries oldest-first.
	sorted := make([]ChangelogEntry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})

	first := sorted[0].Timestamp
	buckets := map[int64]*AggregationWindow{}

	for _, e := range sorted {
		offset := e.Timestamp.Sub(first)
		bucketIdx := int64(offset / windowSize)
		if _, ok := buckets[bucketIdx]; !ok {
			start := first.Add(time.Duration(bucketIdx) * windowSize)
			buckets[bucketIdx] = &AggregationWindow{
				Start: start,
				End:   start.Add(windowSize),
			}
		}
		buckets[bucketIdx].Findings = append(buckets[bucketIdx].Findings, e.Findings...)
	}

	// Collect and sort windows.
	keys := make([]int64, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	result := AggregationResult{}
	var peak *AggregationWindow
	for _, k := range keys {
		w := buckets[k]
		result.Windows = append(result.Windows, *w)
		result.TotalFindings += len(w.Findings)
		if peak == nil || len(w.Findings) > len(peak.Findings) {
			copy := *w
			peak = &copy
		}
	}
	result.PeakWindow = peak
	return result
}
