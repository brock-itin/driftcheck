// Package drift provides drift detection, reporting, and scoring utilities
// for comparing running container state against compose/helm definitions.
//
// # Scoring
//
// The scoring sub-system assigns a numeric weight to each drift finding based
// on its severity level and aggregates per-service totals:
//
//	critical → 10.0
//	high     →  5.0
//	medium   →  2.0
//	low      →  1.0
//	info     →  0.5
//
// Use ScoreFindings to obtain a ranked []DriftScore slice (highest score
// first). Use TotalScore to collapse the slice into a single aggregate value
// suitable for thresholding in CI pipelines.
package drift
