package drift

import (
	"testing"
)

func makeReport(findings ...Finding) Report {
	return Report{Findings: findings}
}

func finding(service, typ string) Finding {
	return Finding{Service: service, Type: typ, Expected: "a", Actual: "b"}
}

func TestFilter_EmptyOptions_ReturnsAll(t *testing.T) {
	r := makeReport(finding("web", "image"), finding("db", "env"))
	out := Filter(r, FilterOptions{})
	if len(out.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(out.Findings))
	}
}

func TestFilter_ByType(t *testing.T) {
	r := makeReport(finding("web", "image"), finding("db", "env"), finding("cache", "image"))
	out := Filter(r, FilterOptions{Types: []string{"image"}})
	if len(out.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(out.Findings))
	}
	for _, f := range out.Findings {
		if f.Type != "image" {
			t.Errorf("unexpected type %q", f.Type)
		}
	}
}

func TestFilter_ByService(t *testing.T) {
	r := makeReport(finding("web", "image"), finding("db", "env"), finding("web", "env"))
	out := Filter(r, FilterOptions{Services: []string{"web"}})
	if len(out.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(out.Findings))
	}
	for _, f := range out.Findings {
		if f.Service != "web" {
			t.Errorf("unexpected service %q", f.Service)
		}
	}
}

func TestFilter_ExcludeServices(t *testing.T) {
	r := makeReport(finding("web", "image"), finding("db", "env"), finding("cache", "image"))
	out := Filter(r, FilterOptions{ExcludeServices: []string{"db"}})
	if len(out.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(out.Findings))
	}
	for _, f := range out.Findings {
		if f.Service == "db" {
			t.Error("excluded service db should not appear")
		}
	}
}

func TestFilter_TypeAndService(t *testing.T) {
	r := makeReport(finding("web", "image"), finding("web", "env"), finding("db", "image"))
	out := Filter(r, FilterOptions{Types: []string{"image"}, Services: []string{"web"}})
	if len(out.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(out.Findings))
	}
	if out.Findings[0].Service != "web" || out.Findings[0].Type != "image" {
		t.Errorf("unexpected finding: %+v", out.Findings[0])
	}
}

func TestFilter_NoMatches(t *testing.T) {
	r := makeReport(finding("web", "image"))
	out := Filter(r, FilterOptions{Services: []string{"db"}})
	if len(out.Findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(out.Findings))
	}
}
