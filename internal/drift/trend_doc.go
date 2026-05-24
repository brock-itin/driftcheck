// Package drift provides drift detection between running containers and their
// source compose/helm definitions.
//
// # Trend Analysis
//
// The trend sub-feature builds a time-ordered series of drift counts from
// saved changelog entries, enabling operators to observe whether configuration
// drift is growing or shrinking over time.
//
// Basic usage:
//
//	entries, err := drift.LoadChangelog("drift-changelog.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	trend := drift.BuildTrend(entries)
//	output.WriteTrend(os.Stdout, trend)
//
// The Trend type exposes:
//   - Points  — ordered slice of TrendPoint (timestamp + totals)
//   - Delta() — change in finding count between the two most recent points
//   - Latest() — the most recent TrendPoint
package drift
