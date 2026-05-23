package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
)

func makeChangelog(entries ...drift.ChangelogEntry) *drift.Changelog {
	return &drift.Changelog{Entries: entries}
}

func TestWriteChangelog_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteChangelog(&buf, makeChangelog()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changelog entries") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteChangelog_WithEntries(t *testing.T) {
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	cl := makeChangelog(
		drift.ChangelogEntry{
			Timestamp:  ts,
			ComposFile: "docker-compose.yml",
			Total:      3,
			BySeverity: map[string]int{"high": 2, "low": 1},
			ByType:     map[string]int{"image": 2, "env": 1},
		},
	)
	var buf bytes.Buffer
	if err := WriteChangelog(&buf, cl); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "docker-compose.yml") {
		t.Errorf("expected compose file in output")
	}
	if !strings.Contains(out, "3") {
		t.Errorf("expected total count in output")
	}
	if !strings.Contains(out, "TIMESTAMP") {
		t.Errorf("expected header row")
	}
}

func TestChangelogSummary_Empty(t *testing.T) {
	s := ChangelogSummary(makeChangelog())
	if !strings.Contains(s, "no previous runs") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestChangelogSummary_WithEntries(t *testing.T) {
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	cl := makeChangelog(
		drift.ChangelogEntry{Timestamp: ts, ComposFile: "a.yml", Total: 1},
		drift.ChangelogEntry{Timestamp: ts.Add(time.Hour), ComposFile: "b.yml", Total: 5},
	)
	s := ChangelogSummary(cl)
	if !strings.Contains(s, "b.yml") {
		t.Errorf("expected last entry compose file in summary: %s", s)
	}
	if !strings.Contains(s, "5") {
		t.Errorf("expected finding count in summary: %s", s)
	}
}
