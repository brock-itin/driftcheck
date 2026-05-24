package drift

import (
	"strings"

	"github.com/yourorg/driftcheck/internal/compose"
)

// LabelDrift represents a mismatch between expected and actual Docker labels.
type LabelDrift struct {
	Service  string
	Label    string
	Expected string
	Actual   string
	Missing  bool
}

// DetectLabelDrift compares labels defined in a compose service against
// labels observed on a running container. It returns a slice of LabelDrift
// entries for every discrepancy found.
func DetectLabelDrift(svc compose.Service, actual map[string]string) []LabelDrift {
	var drifts []LabelDrift

	for k, expected := range svc.Labels {
		actualVal, ok := actual[k]
		if !ok {
			drifts = append(drifts, LabelDrift{
				Service:  svc.Name,
				Label:    k,
				Expected: expected,
				Missing:  true,
			})
			continue
		}
		if !strings.EqualFold(actualVal, expected) {
			drifts = append(drifts, LabelDrift{
				Service:  svc.Name,
				Label:    k,
				Expected: expected,
				Actual:   actualVal,
			})
		}
	}

	return drifts
}

// LabelDriftToFindings converts []LabelDrift into []Finding so they can be
// included in the standard drift report pipeline.
func LabelDriftToFindings(drifts []LabelDrift) []Finding {
	findings := make([]Finding, 0, len(drifts))
	for _, d := range drifts {
		desc := "label value mismatch"
		if d.Missing {
			desc = "label missing on container"
		}
		findings = append(findings, Finding{
			Service: d.Service,
			Type:    "label",
			Field:   d.Label,
			Expected: d.Expected,
			Actual:   d.Actual,
			Message: desc,
		})
	}
	return findings
}
