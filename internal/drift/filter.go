package drift

// FilterOptions controls which findings are included in results.
type FilterOptions struct {
	// Types restricts findings to specific drift types.
	// If empty, all types are included.
	Types []string

	// Services restricts findings to specific service names.
	// If empty, all services are included.
	Services []string

	// ExcludeServices omits findings for the listed service names.
	ExcludeServices []string
}

// Filter returns a new Report containing only findings that match opts.
func Filter(r Report, opts FilterOptions) Report {
	if len(opts.Types) == 0 && len(opts.Services) == 0 && len(opts.ExcludeServices) == 0 {
		return r
	}

	typeSet := toSet(opts.Types)
	serviceSet := toSet(opts.Services)
	excludeSet := toSet(opts.ExcludeServices)

	var filtered []Finding
	for _, f := range r.Findings {
		if len(excludeSet) > 0 && excludeSet[f.Service] {
			continue
		}
		if len(serviceSet) > 0 && !serviceSet[f.Service] {
			continue
		}
		if len(typeSet) > 0 && !typeSet[f.Type] {
			continue
		}
		filtered = append(filtered, f)
	}

	return Report{Findings: filtered}
}

func toSet(items []string) map[string]bool {
	if len(items) == 0 {
		return nil
	}
	s := make(map[string]bool, len(items))
	for _, v := range items {
		s[v] = true
	}
	return s
}
