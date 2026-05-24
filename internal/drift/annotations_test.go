package drift

import (
	"strings"
	"testing"
)

func TestParseAnnotations_Valid(t *testing.T) {
	raw := []string{"env=production", "team=platform"}
	got, err := ParseAnnotations(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(got))
	}
	if got[0].Key != "env" || got[0].Value != "production" {
		t.Errorf("unexpected first annotation: %+v", got[0])
	}
	if got[1].Key != "team" || got[1].Value != "platform" {
		t.Errorf("unexpected second annotation: %+v", got[1])
	}
}

func TestParseAnnotations_MissingEquals(t *testing.T) {
	_, err := ParseAnnotations([]string{"badformat"})
	if err == nil {
		t.Fatal("expected error for missing '=' but got nil")
	}
}

func TestParseAnnotations_EmptyKey(t *testing.T) {
	_, err := ParseAnnotations([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key but got nil")
	}
}

func TestParseAnnotations_Empty(t *testing.T) {
	got, err := ParseAnnotations(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestAnnotateFindings_AllServices(t *testing.T) {
	r := &Report{
		Findings: []Finding{
			{Service: "web", Type: "image", Message: "image mismatch"},
			{Service: "db", Type: "env", Message: "env mismatch"},
		},
	}
	anns := []Annotation{{Key: "env", Value: "staging"}}
	result := AnnotateFindings(r, "", anns)

	for _, f := range result.Findings {
		if !strings.Contains(f.Message, "env=staging") {
			t.Errorf("expected annotation in message, got: %q", f.Message)
		}
	}
}

func TestAnnotateFindings_FilteredService(t *testing.T) {
	r := &Report{
		Findings: []Finding{
			{Service: "web", Type: "image", Message: "image mismatch"},
			{Service: "db", Type: "env", Message: "env mismatch"},
		},
	}
	anns := []Annotation{{Key: "owner", Value: "alice"}}
	result := AnnotateFindings(r, "web", anns)

	if !strings.Contains(result.Findings[0].Message, "owner=alice") {
		t.Errorf("expected annotation on web finding, got: %q", result.Findings[0].Message)
	}
	if strings.Contains(result.Findings[1].Message, "owner=alice") {
		t.Errorf("expected no annotation on db finding, got: %q", result.Findings[1].Message)
	}
}

func TestAnnotateFindings_NilReport(t *testing.T) {
	result := AnnotateFindings(nil, "", []Annotation{{Key: "k", Value: "v"}})
	if result != nil {
		t.Errorf("expected nil result for nil report")
	}
}

func TestAnnotateFindings_NoAnnotations(t *testing.T) {
	r := &Report{
		Findings: []Finding{
			{Service: "web", Type: "image", Message: "original"},
		},
	}
	result := AnnotateFindings(r, "", nil)
	if result.Findings[0].Message != "original" {
		t.Errorf("expected unchanged message, got: %q", result.Findings[0].Message)
	}
}
