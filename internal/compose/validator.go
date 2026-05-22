package compose

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError holds a list of validation issues found in a compose file.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("compose validation failed: %s", strings.Join(e.Issues, "; "))
}

// Validate checks a parsed ComposeFile for common configuration issues
// that could affect drift detection accuracy.
func Validate(cf *ComposeFile) error {
	if cf == nil {
		return errors.New("compose file is nil")
	}

	var issues []string

	if len(cf.Services) == 0 {
		issues = append(issues, "no services defined")
	}

	for name, svc := range cf.Services {
		if name == "" {
			issues = append(issues, "service has empty name")
			continue
		}
		if svc.Image == "" {
			issues = append(issues, fmt.Sprintf("service %q has no image defined", name))
		}
		for _, env := range svc.Environment {
			if !strings.Contains(env, "=") {
				issues = append(issues,
					fmt.Sprintf("service %q has malformed env entry %q (missing '=')", name, env))
			}
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}
