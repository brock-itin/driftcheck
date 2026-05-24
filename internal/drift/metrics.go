package drift

import (
	"fmt"
	"sort"
	"time"
)

// Metrics holds aggregated statistics for a drift detection run.
type Metrics struct {
	RunAt          time.Time
	ServicesTotal  int
	ServicesDrifted int
	FindingsTotal  int
	BySeverity     map[string]int
	ByType         map[string]int
	DurationMs     int64
}

// CollectMetrics builds a Metrics snapshot from a Report and elapsed duration.
func CollectMetrics(r Report, elapsed time.Duration) Metrics {
	m := Metrics{
		RunAt:      time.Now().UTC(),
		BySeverity: make(map[string]int),
		ByType:     make(map[string]int),
		DurationMs: elapsed.Milliseconds(),
	}

	services := make(map[string]struct{})
	drifted := make(map[string]struct{})

	for _, f := range r.Findings {
		m.FindingsTotal++
		m.BySeverity[f.Severity.String()]++
		m.ByType[string(f.Type)]++
		services[f.Service] = struct{}{}
		drifted[f.Service] = struct{}{}
	}

	m.ServicesTotal = len(services)
	m.ServicesDrifted = len(drifted)
	return m
}

// TopDriftedServices returns up to n service names sorted by finding count descending.
func TopDriftedServices(r Report, n int) []string {
	counts := make(map[string]int)
	for _, f := range r.Findings {
		counts[f.Service]++
	}
	type pair struct {
		name  string
		count int
	}
	pairs := make([]pair, 0, len(counts))
	for k, v := range counts {
		pairs = append(pairs, pair{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].name < pairs[j].name
	})
	result := make([]string, 0, n)
	for i, p := range pairs {
		if i >= n {
			break
		}
		result = append(result, fmt.Sprintf("%s (%d)", p.name, p.count))
	}
	return result
}
