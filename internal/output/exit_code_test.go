package output_test

import (
	"testing"

	"github.com/user/driftcheck/internal/drift"
	"github.com/user/driftcheck/internal/output"
)

func TestResolveExitCode_NoFindings(t *testing.T) {
	r := drift.Report{}
	code := output.ResolveExitCode(r)
	if code != output.ExitOK {
		t.Errorf("expected ExitOK (%d), got %d", output.ExitOK, code)
	}
}

func TestResolveExitCode_WithFindings(t *testing.T) {
	r := drift.Report{
		Findings: []drift.Finding{
			{
				Service: "web",
				Field:   "image",
				Expected: "nginx:1.25",
				Actual:   "nginx:latest",
			},
		},
	}
	code := output.ResolveExitCode(r)
	if code != output.ExitDrift {
		t.Errorf("expected ExitDrift (%d), got %d", output.ExitDrift, code)
	}
}

func TestResolveExitCode_MultipleFindings(t *testing.T) {
	r := drift.Report{
		Findings: []drift.Finding{
			{Service: "api", Field: "image", Expected: "go:1.21", Actual: "go:1.20"},
			{Service: "db", Field: "env.DB_PASS", Expected: "secret", Actual: "wrong"},
		},
	}
	code := output.ResolveExitCode(r)
	if code != output.ExitDrift {
		t.Errorf("expected ExitDrift (%d), got %d", output.ExitDrift, code)
	}
}

func TestExitCodeConstants(t *testing.T) {
	if output.ExitOK != 0 {
		t.Errorf("ExitOK should be 0, got %d", output.ExitOK)
	}
	if output.ExitDrift != 1 {
		t.Errorf("ExitDrift should be 1, got %d", output.ExitDrift)
	}
	if output.ExitError != 2 {
		t.Errorf("ExitError should be 2, got %d", output.ExitError)
	}
}
