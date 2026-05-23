package drift

import (
	"context"
	"log"
	"time"

	"github.com/yourorg/driftcheck/internal/compose"
	"github.com/yourorg/driftcheck/internal/snapshot"
)

// WatchOptions configures the drift watcher.
type WatchOptions struct {
	Interval    time.Duration
	ComposeFile string
	RuleSet     RuleSet
	OnDrift     func(Report)
}

// Watcher polls for drift at a fixed interval and invokes a callback.
type Watcher struct {
	opts    WatchOptions
	builder *snapshot.Builder
}

// NewWatcher creates a Watcher with the given options and snapshot builder.
func NewWatcher(opts WatchOptions, builder *snapshot.Builder) *Watcher {
	if opts.Interval <= 0 {
		opts.Interval = 60 * time.Second
	}
	return &Watcher{opts: opts, builder: builder}
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.check(); err != nil {
				log.Printf("driftcheck watcher: check error: %v", err)
			}
		}
	}
}

func (w *Watcher) check() error {
	svc, err := compose.ParseFile(w.opts.ComposeFile)
	if err != nil {
		return err
	}

	snap, err := w.builder.Build(context.Background())
	if err != nil {
		return err
	}

	report := Detect(svc, snap, w.opts.RuleSet)
	if len(report.Findings) > 0 && w.opts.OnDrift != nil {
		w.opts.OnDrift(report)
	}
	return nil
}
