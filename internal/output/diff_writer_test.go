package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftcheck/internal/drift"
	"github.com/yourorg/driftcheck/internal/output"
)

func makeDiffReport(findings ...drift.Finding) drift.Report {
	return drift.Report{Findings: findings}
}

func diffFinding(svc, field, expected, actual string) drift.Finding {
	return drift.Finding{
		Service:  svc,
		Field:    field,
		Expected: expected,
		Actual:   actual,
	}
}

func TestWriteDiff_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	err := output.WriteDiff(&buf, makeDiffReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriteDiff_WithFindings(t *testing.T) {
	r := makeDiffReport(
		diffFinding("web", "image", "nginx:1.24", "nginx:1.21"),
		diffFinding("web", "env.PORT", "8080", "9090"),
	)
	var buf bytes.Buffer
	err := output.WriteDiff(&buf, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "--- expected/web") {
		t.Errorf("expected expected header, got: %s", out)
	}
	if !strings.Contains(out, "+++ running/web") {
		t.Errorf("expected running header, got: %s", out)
	}
	if !strings.Contains(out, "- nginx:1.24") {
		t.Errorf("expected old image line, got: %s", out)
	}
	if !strings.Contains(out, "+ nginx:1.21") {
		t.Errorf("expected new image line, got: %s", out)
	}
	if !strings.Contains(out, "@@ image @@") {
		t.Errorf("expected field hunk header, got: %s", out)
	}
}

func TestWriteDiff_MultipleServices_SortedOutput(t *testing.T) {
	r := makeDiffReport(
		diffFinding("zebra", "image", "old", "new"),
		diffFinding("alpha", "image", "old", "new"),
	)
	var buf bytes.Buffer
	_ = output.WriteDiff(&buf, r)
	out := buf.String()
	alphaIdx := strings.Index(out, "alpha")
	zebraIdx := strings.Index(out, "zebra")
	if alphaIdx > zebraIdx {
		t.Errorf("expected alpha before zebra in output")
	}
}

func TestWriteDiff_EmptyExpected(t *testing.T) {
	r := makeDiffReport(diffFinding("db", "env.SECRET", "", "hunter2"))
	var buf bytes.Buffer
	_ = output.WriteDiff(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "- (not set)") {
		t.Errorf("expected '(not set)' for empty expected, got: %s", out)
	}
}
