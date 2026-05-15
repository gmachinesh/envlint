package sorter_test

import (
	"testing"

	"github.com/nicholasgasior/envlint/sorter"
	"github.com/nicholasgasior/envlint/validator"
)

func makeResults(pairs ...string) []validator.Result {
	var out []validator.Result
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, validator.Result{Key: pairs[i], Severity: pairs[i+1]})
	}
	return out
}

func TestSortKeys_Asc(t *testing.T) {
	keys := []string{"ZEBRA", "apple", "Mango"}
	sorted := sorter.SortKeys(keys, sorter.Options{Order: sorter.OrderAsc})
	expected := []string{"apple", "Mango", "ZEBRA"}
	for i, k := range sorted {
		if k != expected[i] {
			t.Errorf("pos %d: got %q want %q", i, k, expected[i])
		}
	}
}

func TestSortKeys_Desc(t *testing.T) {
	keys := []string{"apple", "ZEBRA", "Mango"}
	sorted := sorter.SortKeys(keys, sorter.Options{Order: sorter.OrderDesc})
	expected := []string{"ZEBRA", "Mango", "apple"}
	for i, k := range sorted {
		if k != expected[i] {
			t.Errorf("pos %d: got %q want %q", i, k, expected[i])
		}
	}
}

func TestSortKeys_DoesNotMutateOriginal(t *testing.T) {
	orig := []string{"C", "A", "B"}
	_ = sorter.SortKeys(orig, sorter.DefaultOptions())
	if orig[0] != "C" {
		t.Error("original slice was mutated")
	}
}

func TestSortResults_Asc(t *testing.T) {
	results := makeResults("ZEBRA", "ok", "APPLE", "error", "MANGO", "warning")
	sorted := sorter.SortResults(results, sorter.Options{Order: sorter.OrderAsc})
	if sorted[0].Key != "APPLE" || sorted[1].Key != "MANGO" || sorted[2].Key != "ZEBRA" {
		t.Errorf("unexpected order: %v", sorted)
	}
}

func TestSortResults_Desc(t *testing.T) {
	results := makeResults("APPLE", "ok", "ZEBRA", "error", "MANGO", "warning")
	sorted := sorter.SortResults(results, sorter.Options{Order: sorter.OrderDesc})
	if sorted[0].Key != "ZEBRA" || sorted[1].Key != "MANGO" || sorted[2].Key != "APPLE" {
		t.Errorf("unexpected order: %v", sorted)
	}
}

func TestSortResults_Severity(t *testing.T) {
	results := makeResults("B", "info", "A", "error", "C", "warning", "D", "ok")
	sorted := sorter.SortResults(results, sorter.Options{Order: sorter.OrderSeverity})
	expected := []string{"A", "C", "B", "D"}
	for i, r := range sorted {
		if r.Key != expected[i] {
			t.Errorf("pos %d: got %q want %q", i, r.Key, expected[i])
		}
	}
}

func TestSortResults_SeverityTieBreakByKey(t *testing.T) {
	results := makeResults("Z", "error", "A", "error")
	sorted := sorter.SortResults(results, sorter.Options{Order: sorter.OrderSeverity})
	if sorted[0].Key != "A" {
		t.Errorf("expected A first, got %q", sorted[0].Key)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := sorter.DefaultOptions()
	if opts.Order != sorter.OrderAsc {
		t.Errorf("expected OrderAsc, got %q", opts.Order)
	}
}
