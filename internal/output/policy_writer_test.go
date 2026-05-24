package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftcheck/internal/drift"
)

func makePolicyFindings() []drift.Finding {
	return []drift.Finding{
		{Service: "api", Type: "image", Detail: "tag mismatch: latest vs v1.2"},
		{Service: "web", Type: "env", Detail: "PORT missing"},
	}
}

func TestWritePolicyEvaluation_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	policy := &drift.Policy{}
	if err := WritePolicyEvaluation(&buf, nil, policy); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "skipped") {
		t.Errorf("expected skip message, got: %s", buf.String())
	}
}

func TestWritePolicyEvaluation_WithFindings(t *testing.T) {
	var buf bytes.Buffer
	policy := &drift.Policy{
		Rules: []drift.PolicyRule{
			{Type: "image", Action: drift.PolicyActionFail},
		},
	}
	if err := WritePolicyEvaluation(&buf, makePolicyFindings(), policy); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "fail") {
		t.Errorf("expected 'fail' in output, got: %s", out)
	}
	if !strings.Contains(out, "warn") {
		t.Errorf("expected 'warn' in output, got: %s", out)
	}
	if !strings.Contains(out, "api") {
		t.Errorf("expected service 'api' in output, got: %s", out)
	}
}

func TestPolicyActionCounts_Empty(t *testing.T) {
	policy := &drift.Policy{}
	counts := PolicyActionCounts(nil, policy)
	if len(counts) != 0 {
		t.Errorf("expected empty counts, got %v", counts)
	}
}

func TestPolicyActionCounts_Mixed(t *testing.T) {
	policy := &drift.Policy{
		Rules: []drift.PolicyRule{
			{Type: "image", Action: drift.PolicyActionFail},
			{Type: "env", Action: drift.PolicyActionIgnore},
		},
	}
	counts := PolicyActionCounts(makePolicyFindings(), policy)
	if counts[drift.PolicyActionFail] != 1 {
		t.Errorf("expected 1 fail, got %d", counts[drift.PolicyActionFail])
	}
	if counts[drift.PolicyActionIgnore] != 1 {
		t.Errorf("expected 1 ignore, got %d", counts[drift.PolicyActionIgnore])
	}
}
