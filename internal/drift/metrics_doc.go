// Package drift provides drift detection between running containers and their
// source compose/helm definitions.
//
// # Metrics
//
// The Metrics type captures aggregated statistics from a single detection run,
// including total findings, per-severity and per-type breakdowns, and elapsed
// wall-clock duration.
//
// Usage:
//
//	start := time.Now()
//	report := drift.Detect(compose, containers)
//	metrics := drift.CollectMetrics(report, time.Since(start))
//	output.WriteMetrics(os.Stdout, metrics)
//
// TopDriftedServices returns the N most-affected service names in descending
// order of finding count, useful for dashboard summaries and alerting.
package drift
