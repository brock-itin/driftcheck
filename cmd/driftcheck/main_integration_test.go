//go:build integration

package main

import (
	"os/exec"
	"strings"
	"testing"
)

// TestCLI_Help verifies the binary responds to --help without error.
func TestCLI_Help(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "--help")
	out, err := cmd.CombinedOutput()
	// flag package exits with code 2 for --help; that's expected
	if err != nil && !strings.Contains(err.Error(), "exit status 2") {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	for _, want := range []string{"file", "format", "exit-code"} {
		if !strings.Contains(string(out), want) {
			t.Errorf("expected flag %q in help output", want)
		}
	}
}

// TestCLI_MissingFile verifies a missing compose file produces a clear error.
func TestCLI_MissingFile(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "--file", "/nonexistent/docker-compose.yml")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for missing file")
	}
	if !strings.Contains(string(out), "error") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}
