package comparator_test

import (
	"testing"

	"github.com/nicholasgasior/envlint/comparator"
)

func TestCompare_NoDiff(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	head := map[string]string{"A": "1", "B": "2"}
	r := comparator.Compare(base, head)
	if r.HasDiff() {
		t.Fatal("expected no diff")
	}
}

func TestCompare_Added(t *testing.T) {
	base := map[string]string{"A": "1"}
	head := map[string]string{"A": "1", "B": "2"}
	r := comparator.Compare(base, head)
	if !r.HasDiff() {
		t.Fatal("expected diff")
	}
	if r.Changes[1].Kind != comparator.Added || r.Changes[1].Key != "B" {
		t.Errorf("unexpected change: %+v", r.Changes)
	}
}

func TestCompare_Removed(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	head := map[string]string{"A": "1"}
	r := comparator.Compare(base, head)
	found := false
	for _, c := range r.Changes {
		if c.Key == "B" && c.Kind == comparator.Removed {
			found = true
		}
	}
	if !found {
		t.Error("expected B to be removed")
	}
}

func TestCompare_Modified(t *testing.T) {
	base := map[string]string{"A": "old"}
	head := map[string]string{"A": "new"}
	r := comparator.Compare(base, head)
	if len(r.Changes) != 1 || r.Changes[0].Kind != comparator.Modified {
		t.Fatalf("expected modified, got %+v", r.Changes)
	}
	if r.Changes[0].OldValue != "old" || r.Changes[0].NewValue != "new" {
		t.Errorf("unexpected values: %+v", r.Changes[0])
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	base := map[string]string{"Z": "1", "A": "2"}
	head := map[string]string{"Z": "1", "A": "2"}
	r := comparator.Compare(base, head)
	if r.Changes[0].Key != "A" {
		t.Errorf("expected A first, got %s", r.Changes[0].Key)
	}
}
