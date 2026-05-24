package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/driftcheck/internal/drift"
)

func makeMetrics(findings int, services int, drifted int) drift.Metrics {
	m := drift.Metrics{
		RunAt:           time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		ServicesTotal:   services,
		ServicesDrifted: drifted,
		FindingsTotal:   findings,
		DurationMs:      42,
		BySeverity:      map[string]int{},
		ByType:          map[string]int{},
	}
	if findings > 0 {
		m.BySeverity["high"] = findings
		m.ByType["image"] = findings
	}
	return m
}

func TestWriteMetrics_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	WriteMetrics(&buf, makeMetrics(0, 3, 0))
	out := buf.String()
	if !strings.Contains(out, "Services scanned:") {
		t.Error("expected 'Services scanned:' in output")
	}
	if !strings.Contains(out, "Total findings:  0") {
		t.Errorf("expected zero findings line, got:\n%s", out)
	}
	if strings.Contains(out, "by severity") {
		t.Error("should not show severity section when no findings")
	}
}

func TestWriteMetrics_WithFindings(t *testing.T) {
	var buf bytes.Buffer
	WriteMetrics(&buf, makeMetrics(4, 2, 2))
	out := buf.String()
	if !strings.Contains(out, "high") {
		t.Error("expected severity 'high' in output")
	}
	if !strings.Contains(out, "image") {
		t.Error("expected type 'image' in output")
	}
	if !strings.Contains(out, "42ms") {
		t.Error("expected duration in output")
	}
}

func TestMetricsSummaryLine(t *testing.T) {
	m := makeMetrics(3, 5, 2)
	line := MetricsSummaryLine(m)
	if !strings.Contains(line, "5 service(s)") {
		t.Errorf("expected total services in summary, got: %s", line)
	}
	if !strings.Contains(line, "3 finding(s)") {
		t.Errorf("expected finding count in summary, got: %s", line)
	}
	if !strings.Contains(line, "42ms") {
		t.Errorf("expected duration in summary, got: %s", line)
	}
}

func TestSortedMapKeys_Order(t *testing.T) {
	m := map[string]int{"zebra": 1, "apple": 2, "mango": 3}
	keys := sortedMapKeys(m)
	if keys[0] != "apple" || keys[1] != "mango" || keys[2] != "zebra" {
		t.Errorf("unexpected order: %v", keys)
	}
}
