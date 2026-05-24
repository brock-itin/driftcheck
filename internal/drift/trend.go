package drift

import (
	"sort"
	"time"
)

// TrendPoint represents the drift count at a specific point in time.
type TrendPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Total     int       `json:"total"`
	ByType    map[string]int `json:"by_type"`
}

// Trend holds an ordered series of TrendPoints derived from changelog entries.
type Trend struct {
	Points []TrendPoint `json:"points"`
}

// BuildTrend constructs a Trend from a slice of ChangelogEntry values.
func BuildTrend(entries []ChangelogEntry) Trend {
	points := make([]TrendPoint, 0, len(entries))
	for _, e := range entries {
		byType := make(map[string]int)
		for _, f := range e.Findings {
			byType[string(f.Type)]++
		}
		points = append(points, TrendPoint{
			Timestamp: e.RecordedAt,
			Total:     len(e.Findings),
			ByType:    byType,
		})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Timestamp.Before(points[j].Timestamp)
	})
	return Trend{Points: points}
}

// Delta returns the difference in total drift count between the last two points.
// Returns 0 if fewer than two points exist.
func (t Trend) Delta() int {
	if len(t.Points) < 2 {
		return 0
	}
	last := t.Points[len(t.Points)-1]
	prev := t.Points[len(t.Points)-2]
	return last.Total - prev.Total
}

// Latest returns the most recent TrendPoint, or a zero value if empty.
func (t Trend) Latest() TrendPoint {
	if len(t.Points) == 0 {
		return TrendPoint{}
	}
	return t.Points[len(t.Points)-1]
}
