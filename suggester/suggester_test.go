package suggester_test

import (
	"strings"
	"testing"

	"github.com/user/envlint/suggester"
	"github.com/user/envlint/validator"
)

func makeResult(key, message string, ok bool) validator.Result {
	return validator.Result{
		Key:     key,
		Message: message,
		OK:      ok,
	}
}

func TestSuggest_NoResults(t *testing.T) {
	got := suggester.Suggest(nil)
	if len(got) != 0 {
		t.Fatalf("expected 0 suggestions, got %d", len(got))
	}
}

func TestSuggest_SkipsOKResults(t *testing.T) {
	results := []validator.Result{
		makeResult("PORT", "ok", true),
		makeResult("HOST", "ok", true),
	}
	got := suggester.Suggest(results)
	if len(got) != 0 {
		t.Fatalf("expected 0 suggestions for passing results, got %d", len(got))
	}
}

func TestSuggest_MissingVariable(t *testing.T) {
	results := []validator.Result{
		makeResult("DATABASE_URL", "missing required variable", false),
	}
	got := suggester.Suggest(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if got[0].Key != "DATABASE_URL" {
		t.Errorf("expected key DATABASE_URL, got %q", got[0].Key)
	}
	if !strings.Contains(got[0].Message, "Add") {
		t.Errorf("expected 'Add' hint, got: %q", got[0].Message)
	}
}

func TestSuggest_PatternMismatch(t *testing.T) {
	results := []validator.Result{
		makeResult("API_URL", "value does not match pattern", false),
	}
	got := suggester.Suggest(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Message, "pattern") {
		t.Errorf("expected pattern hint, got: %q", got[0].Message)
	}
}

func TestSuggest_EmptyValue(t *testing.T) {
	results := []validator.Result{
		makeResult("SECRET_KEY", "value is empty", false),
	}
	got := suggester.Suggest(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Message, "empty") {
		t.Errorf("expected empty hint, got: %q", got[0].Message)
	}
}

func TestSuggest_DefaultHint(t *testing.T) {
	results := []validator.Result{
		makeResult("UNKNOWN_VAR", "some unknown error", false),
	}
	got := suggester.Suggest(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Message, "schema") {
		t.Errorf("expected schema hint, got: %q", got[0].Message)
	}
}
