package filter_test

import (
	"testing"

	"github.com/yourorg/envlint/filter"
	"github.com/yourorg/envlint/validator"
)

func makeResults() []validator.Result {
	return []validator.Result{
		{Variable: "APP_HOST", Severity: "error", Message: "missing required variable"},
		{Variable: "APP_PORT", Severity: "warning", Message: "pattern mismatch"},
		{Variable: "DB_URL", Severity: "error", Message: "missing required variable"},
		{Variable: "LOG_LEVEL", Severity: "info", Message: "optional variable not set"},
	}
}

func TestApply_NoFilter(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{})
	if len(got) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(got))
	}
}

func TestApply_SeverityError(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{Severity: "error"})
	if len(got) != 2 {
		t.Errorf("expected 2 error results, got %d", len(got))
	}
	for _, r := range got {
		if r.Severity != "error" {
			t.Errorf("expected severity 'error', got %q", r.Severity)
		}
	}
}

func TestApply_SeverityWarning(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{Severity: "warning"})
	// warning and error should be included
	if len(got) != 3 {
		t.Errorf("expected 3 results (warning+error), got %d", len(got))
	}
}

func TestApply_SeverityInfo(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{Severity: "info"})
	if len(got) != 4 {
		t.Errorf("expected all 4 results for info level, got %d", len(got))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{Prefix: "APP_"})
	if len(got) != 2 {
		t.Errorf("expected 2 APP_ results, got %d", len(got))
	}
	for _, r := range got {
		if r.Variable[:4] != "APP_" {
			t.Errorf("expected variable with prefix APP_, got %q", r.Variable)
		}
	}
}

func TestApply_CombinedFilter(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{Severity: "error", Prefix: "APP_"})
	if len(got) != 1 {
		t.Errorf("expected 1 result, got %d", len(got))
	}
	if got[0].Variable != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %q", got[0].Variable)
	}
}

func TestApply_NoMatches(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{Prefix: "NONEXISTENT_"})
	if len(got) != 0 {
		t.Errorf("expected 0 results, got %d", len(got))
	}
}
