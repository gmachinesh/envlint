package grouper_test

import (
	"testing"

	"github.com/yourorg/envlint/grouper"
	"github.com/yourorg/envlint/validator"
)

func makeResults() []validator.Result {
	return []validator.Result{
		{Key: "DB_HOST", OK: true, Severity: "info", Message: ""},
		{Key: "DB_PASS", OK: false, Severity: "error", Message: "missing required"},
		{Key: "APP_ENV", OK: true, Severity: "info", Message: ""},
		{Key: "APP_PORT", OK: false, Severity: "warning", Message: "pattern mismatch"},
		{Key: "SECRET", OK: false, Severity: "error", Message: "missing required"},
	}
}

func TestApply_ByPrefix(t *testing.T) {
	groups := grouper.Apply(makeResults(), grouper.ByPrefix)
	names := map[string]bool{}
	for _, g := range groups {
		names[g.Name] = true
	}
	if !names["app"] || !names["db"] || !names["(none)"] {
		t.Errorf("expected groups app, db, (none); got %v", names)
	}
}

func TestApply_BySeverity(t *testing.T) {
	groups := grouper.Apply(makeResults(), grouper.BySeverity)
	names := map[string]int{}
	for _, g := range groups {
		names[g.Name] = len(g.Results)
	}
	if names["error"] != 2 {
		t.Errorf("expected 2 error results, got %d", names["error"])
	}
	if names["warning"] != 1 {
		t.Errorf("expected 1 warning result, got %d", names["warning"])
	}
	if names["info"] != 2 {
		t.Errorf("expected 2 info results, got %d", names["info"])
	}
}

func TestApply_ByStatus(t *testing.T) {
	groups := grouper.Apply(makeResults(), grouper.ByStatus)
	counts := map[string]int{}
	for _, g := range groups {
		counts[g.Name] = len(g.Results)
	}
	if counts["ok"] != 2 {
		t.Errorf("expected 2 ok, got %d", counts["ok"])
	}
	if counts["fail"] != 3 {
		t.Errorf("expected 3 fail, got %d", counts["fail"])
	}
}

func TestApply_EmptyResults(t *testing.T) {
	groups := grouper.Apply([]validator.Result{}, grouper.BySeverity)
	if len(groups) != 0 {
		t.Errorf("expected no groups for empty input, got %d", len(groups))
	}
}

func TestApply_SortedOutput(t *testing.T) {
	groups := grouper.Apply(makeResults(), grouper.ByPrefix)
	for i := 1; i < len(groups); i++ {
		if groups[i-1].Name > groups[i].Name {
			t.Errorf("groups not sorted: %s > %s", groups[i-1].Name, groups[i].Name)
		}
	}
}
