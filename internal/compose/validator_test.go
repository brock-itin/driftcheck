package compose

import (
	"strings"
	"testing"
)

func makeCompose(services map[string]Service) *ComposeFile {
	return &ComposeFile{Services: services}
}

func TestValidate_NilCompose(t *testing.T) {
	err := Validate(nil)
	if err == nil {
		t.Fatal("expected error for nil compose file")
	}
}

func TestValidate_NoServices(t *testing.T) {
	cf := makeCompose(map[string]Service{})
	err := Validate(cf)
	if err == nil {
		t.Fatal("expected error when no services defined")
	}
	if !strings.Contains(err.Error(), "no services defined") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_MissingImage(t *testing.T) {
	cf := makeCompose(map[string]Service{
		"web": {Image: "", Environment: []string{"PORT=8080"}},
	})
	err := Validate(cf)
	if err == nil {
		t.Fatal("expected error for service with no image")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) != 1 || !strings.Contains(ve.Issues[0], "no image") {
		t.Errorf("unexpected issues: %v", ve.Issues)
	}
}

func TestValidate_MalformedEnv(t *testing.T) {
	cf := makeCompose(map[string]Service{
		"api": {Image: "myapp:latest", Environment: []string{"GOOD=value", "BADENTRY"}},
	})
	err := Validate(cf)
	if err == nil {
		t.Fatal("expected error for malformed env entry")
	}
	if !strings.Contains(err.Error(), "BADENTRY") {
		t.Errorf("error should mention the bad entry, got: %v", err)
	}
}

func TestValidate_Valid(t *testing.T) {
	cf := makeCompose(map[string]Service{
		"web": {Image: "nginx:1.25", Environment: []string{"PORT=80"}},
		"db": {Image: "postgres:15", Environment: []string{"POSTGRES_DB=app"}},
	})
	if err := Validate(cf); err != nil {
		t.Errorf("expected no error for valid compose, got: %v", err)
	}
}

func TestValidationError_Message(t *testing.T) {
	ve := &ValidationError{Issues: []string{"issue one", "issue two"}}
	msg := ve.Error()
	if !strings.Contains(msg, "issue one") || !strings.Contains(msg, "issue two") {
		t.Errorf("error message missing issues: %s", msg)
	}
}
