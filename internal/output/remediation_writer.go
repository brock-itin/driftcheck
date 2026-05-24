package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/driftcheck/internal/drift"
)

// WriteRemediation writes a human-readable remediation plan to w.
func WriteRemediation(w io.Writer, actions []drift.RemediationAction) error {
	if len(actions) == 0 {
		_, err := fmt.Fprintln(w, "No remediation actions required.")
		return err
	}

	_, err := fmt.Fprintf(w, "Remediation Plan (%d action(s)):\n", len(actions))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, strings.Repeat("-", 50))
	if err != nil {
		return err
	}

	for i, a := range actions {
		_, err = fmt.Fprintf(w, "%d. [%s] %s\n", i+1, a.Service, a.Description)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "   $ %s\n\n", a.Command)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemediationSummaryLine returns a one-line summary of available actions.
func RemediationSummaryLine(actions []drift.RemediationAction) string {
	if len(actions) == 0 {
		return "remediation: nothing to do"
	}
	services := make(map[string]struct{})
	for _, a := range actions {
		services[a.Service] = struct{}{}
	}
	return fmt.Sprintf("remediation: %d action(s) across %d service(s)", len(actions), len(services))
}
