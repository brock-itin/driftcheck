package output

import (
	"strings"
	"testing"

	"github.com/user/driftcheck/internal/drift"
)

func makeActions(pairs ...string) []drift.RemediationAction {
	actions := make([]drift.RemediationAction, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		actions = append(actions, drift.RemediationAction{
			Service:     pairs[i],
			DriftType:   "image",
			Description: "Re-deploy " + pairs[i],
			Command:     "docker compose up -d " + pairs[i+1],
		})
	}
	return actions
}

func TestWriteRemediation_NoActions(t *testing.T) {
	var sb strings.Builder
	if err := WriteRemediation(&sb, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sb.String(), "No remediation") {
		t.Errorf("expected no-action message, got: %s", sb.String())
	}
}

func TestWriteRemediation_WithActions(t *testing.T) {
	var sb strings.Builder
	actions := makeActions("web", "web", "api", "api")
	if err := WriteRemediation(&sb, actions); err != nil {
		t.Fatal(err)
	}
	out := sb.String()
	if !strings.Contains(out, "2 action(s)") {
		t.Errorf("expected action count in output: %s", out)
	}
	if !strings.Contains(out, "[web]") {
		t.Errorf("expected service name in output: %s", out)
	}
	if !strings.Contains(out, "$ docker compose") {
		t.Errorf("expected command prefix in output: %s", out)
	}
}

func TestWriteRemediation_NumbersActions(t *testing.T) {
	var sb strings.Builder
	actions := makeActions("svc", "svc")
	_ = WriteRemediation(&sb, actions)
	if !strings.Contains(sb.String(), "1.") {
		t.Errorf("expected numbered action: %s", sb.String())
	}
}

func TestRemediationSummaryLine_Empty(t *testing.T) {
	line := RemediationSummaryLine(nil)
	if !strings.Contains(line, "nothing to do") {
		t.Errorf("unexpected summary: %s", line)
	}
}

func TestRemediationSummaryLine_WithActions(t *testing.T) {
	actions := makeActions("web", "web", "web", "web2", "api", "api")
	line := RemediationSummaryLine(actions)
	if !strings.Contains(line, "3 action(s)") {
		t.Errorf("expected 3 actions in summary: %s", line)
	}
	if !strings.Contains(line, "2 service(s)") {
		t.Errorf("expected 2 services in summary: %s", line)
	}
}
