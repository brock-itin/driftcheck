package compose

import (
	"reflect"
	"testing"
)

func TestServiceNames_Nil(t *testing.T) {
	var c *Compose
	if got := c.ServiceNames(); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestServiceNames_Empty(t *testing.T) {
	c := &Compose{Services: map[string]Service{}}
	if got := c.ServiceNames(); got != nil {
		t.Fatalf("expected nil for empty services, got %v", got)
	}
}

func TestServiceNames_Sorted(t *testing.T) {
	c := &Compose{
		Services: map[string]Service{
			"zebra": {Name: "zebra"},
			"alpha": {Name: "alpha"},
			"mongo": {Name: "mongo"},
		},
	}
	want := []string{"alpha", "mongo", "zebra"}
	got := c.ServiceNames()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestService_LabelsField(t *testing.T) {
	svc := Service{
		Name:   "web",
		Image:  "nginx:latest",
		Labels: map[string]string{"env": "prod", "team": "platform"},
	}
	if svc.Labels["env"] != "prod" {
		t.Errorf("expected label env=prod, got %q", svc.Labels["env"])
	}
}

func TestSortStrings_AlreadySorted(t *testing.T) {
	ss := []string{"a", "b", "c"}
	sortStrings(ss)
	if !reflect.DeepEqual(ss, []string{"a", "b", "c"}) {
		t.Errorf("unexpected result: %v", ss)
	}
}

func TestSortStrings_Reverse(t *testing.T) {
	ss := []string{"z", "m", "a"}
	sortStrings(ss)
	if !reflect.DeepEqual(ss, []string{"a", "m", "z"}) {
		t.Errorf("unexpected result: %v", ss)
	}
}

func TestSortStrings_Single(t *testing.T) {
	ss := []string{"only"}
	sortStrings(ss)
	if ss[0] != "only" {
		t.Errorf("unexpected result: %v", ss)
	}
}
