package drift

import (
	"sort"
	"strings"
)

// NormalizeOptions controls how findings are normalized before comparison.
type NormalizeOptions struct {
	// TrimEnvWhitespace strips leading/trailing whitespace from env values.
	TrimEnvWhitespace bool
	// LowercaseImageRefs normalizes image references to lowercase.
	LowercaseImageRefs bool
	// SortEnvKeys ensures env var keys are compared in a canonical order.
	SortEnvKeys bool
}

// DefaultNormalizeOptions returns sensible normalization defaults.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		TrimEnvWhitespace:  true,
		LowercaseImageRefs: true,
		SortEnvKeys:        true,
	}
}

// NormalizeFindings applies normalization rules to a slice of findings,
// returning a new slice with normalized field values.
func NormalizeFindings(findings []Finding, opts NormalizeOptions) []Finding {
	result := make([]Finding, len(findings))
	for i, f := range findings {
		result[i] = normalizeOne(f, opts)
	}
	return result
}

func normalizeOne(f Finding, opts NormalizeOptions) Finding {
	switch f.Type {
	case "image":
		if opts.LowercaseImageRefs {
			f.Expected = strings.ToLower(f.Expected)
			f.Actual = strings.ToLower(f.Actual)
		}
	case "env":
		if opts.TrimEnvWhitespace {
			f.Expected = strings.TrimSpace(f.Expected)
			f.Actual = strings.TrimSpace(f.Actual)
		}
	}
	return f
}

// NormalizeEnvMap trims whitespace from all values in an env map and
// optionally returns sorted keys for deterministic iteration.
func NormalizeEnvMap(env map[string]string, trim bool) (map[string]string, []string) {
	out := make(map[string]string, len(env))
	keys := make([]string, 0, len(env))
	for k, v := range env {
		if trim {
			v = strings.TrimSpace(v)
		}
		out[k] = v
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return out, keys
}
