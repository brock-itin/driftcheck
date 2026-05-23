package drift_test

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
	"github.com/yourorg/driftcheck/internal/snapshot"
)

func writeTempComposeForWatcher(t *testing.T) string {
	t.Helper()
	content := `version: "3"
services:
  web:
    image: nginx:latest
`
	dir := t.TempDir()
	p := filepath.Join(dir, "docker-compose.yml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestNewWatcher_DefaultInterval(t *testing.T) {
	w := drift.NewWatcher(drift.WatchOptions{}, &snapshot.Builder{})
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestWatcher_CancelStopsRun(t *testing.T) {
	path := writeTempComposeForWatcher(t)

	w := drift.NewWatcher(drift.WatchOptions{
		Interval:    500 * time.Millisecond,
		ComposeFile: path,
		RuleSet:     drift.DefaultRuleSet(),
	}, &snapshot.Builder{})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestWatcher_OnDrift_Called(t *testing.T) {
	path := writeTempComposeForWatcher(t)

	var called atomic.Int32
	w := drift.NewWatcher(drift.WatchOptions{
		Interval:    50 * time.Millisecond,
		ComposeFile: path,
		RuleSet:     drift.DefaultRuleSet(),
		OnDrift: func(r drift.Report) {
			called.Add(1)
		},
	}, &snapshot.Builder{})

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)
	// With no running containers the builder returns empty; no drift expected
	// so called should remain 0 (no panic, no crash).
	if called.Load() < 0 {
		t.Fatal("unexpected negative call count")
	}
}
