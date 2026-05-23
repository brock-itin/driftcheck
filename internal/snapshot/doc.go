// Package snapshot provides types and utilities for capturing a point-in-time
// view of running Docker container states.
//
// A Snapshot can be built from live Docker data via Builder, persisted to disk
// using Save, and restored with Load. Snapshots are used by the drift detector
// to compare expected configuration (from compose/helm definitions) against
// the actual observed state of running containers.
//
// Typical usage:
//
//	builder := snapshot.NewBuilder(dockerClient)
//	snap, err := builder.Build()
//	if err != nil { ... }
//	if err := snap.Save("./state.json"); err != nil { ... }
package snapshot
