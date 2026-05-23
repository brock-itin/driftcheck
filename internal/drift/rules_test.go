package drift

import (
	"testing"
)

func TestDefaultRuleSet_AllEnabled(t *testing.T) {
	rs := DefaultRuleSet()
	if !rs.CheckImage {
		t.Error("expected CheckImage to be true")
	}
	if !rs.CheckEnv {
		t.Error("expected CheckEnv to be true")
	}
	if !rs.CheckPorts {
		t.Error("expected CheckPorts to be true")
	}
}

func TestRuleSetFromFlags_NoFlagsSet_ReturnsDefault(t *testing.T) {
	rs := RuleSetFromFlags(false, false, false, false)
	if !rs.CheckImage || !rs.CheckEnv || !rs.CheckPorts {
		t.Error("expected all checks enabled when no flags set")
	}
}

func TestRuleSetFromFlags_SelectiveFlags(t *testing.T) {
	rs := RuleSetFromFlags(true, false, false, true)
	if !rs.CheckImage {
		t.Error("expected CheckImage true")
	}
	if rs.CheckEnv {
		t.Error("expected CheckEnv false")
	}
	if rs.CheckPorts {
		t.Error("expected CheckPorts false")
	}
}

func TestRuleSet_Enabled(t *testing.T) {
	empty := RuleSet{}
	if empty.Enabled() {
		t.Error("expected Enabled() false for empty RuleSet")
	}

	partial := RuleSet{CheckEnv: true}
	if !partial.Enabled() {
		t.Error("expected Enabled() true when at least one check set")
	}
}

func TestRuleSet_ActiveChecks(t *testing.T) {
	rs := RuleSet{CheckImage: true, CheckPorts: true}
	checks := rs.ActiveChecks()
	if len(checks) != 2 {
		t.Fatalf("expected 2 active checks, got %d", len(checks))
	}

	typeSet := map[FindingType]bool{}
	for _, c := range checks {
		typeSet[c] = true
	}
	if !typeSet[FindingTypeImage] {
		t.Error("expected FindingTypeImage in active checks")
	}
	if !typeSet[FindingTypePorts] {
		t.Error("expected FindingTypePorts in active checks")
	}
	if typeSet[FindingTypeEnv] {
		t.Error("did not expect FindingTypeEnv in active checks")
	}
}

func TestRuleSet_ActiveChecks_AllEnabled(t *testing.T) {
	rs := DefaultRuleSet()
	checks := rs.ActiveChecks()
	if len(checks) != 3 {
		t.Fatalf("expected 3 active checks, got %d", len(checks))
	}
}
