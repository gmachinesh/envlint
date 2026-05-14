package validator_test

import (
	"testing"

	"github.com/user/envlint/schema"
	"github.com/user/envlint/validator"
)

func makeSchema(vars []schema.Variable) *schema.Schema {
	return &schema.Schema{Variables: vars}
}

func TestValidate_AllPresent(t *testing.T) {
	s := makeSchema([]schema.Variable{
		{Name: "APP_ENV", Required: true},
		{Name: "PORT", Required: true},
	})
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}

	report, err := validator.Validate(env, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasErrors() {
		t.Errorf("expected no errors, got failures")
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	s := makeSchema([]schema.Variable{
		{Name: "SECRET_KEY", Required: true},
	})
	env := map[string]string{}

	report, err := validator.Validate(env, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.HasErrors() {
		t.Error("expected errors for missing required variable")
	}
	if len(report.Results) != 1 || report.Results[0].Key != "SECRET_KEY" {
		t.Errorf("unexpected results: %+v", report.Results)
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	s := makeSchema([]schema.Variable{
		{Name: "PORT", Required: true, Pattern: `^\d+$`},
	})
	env := map[string]string{"PORT": "8080"}

	report, err := validator.Validate(env, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasErrors() {
		t.Errorf("expected no errors, got: %+v", report.Results)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	s := makeSchema([]schema.Variable{
		{Name: "PORT", Required: true, Pattern: `^\d+$`},
	})
	env := map[string]string{"PORT": "not-a-number"}

	report, err := validator.Validate(env, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.HasErrors() {
		t.Error("expected pattern mismatch error")
	}
}

func TestValidate_NilSchema(t *testing.T) {
	_, err := validator.Validate(map[string]string{}, nil)
	if err == nil {
		t.Error("expected error for nil schema")
	}
}

func TestValidate_OptionalMissing(t *testing.T) {
	s := makeSchema([]schema.Variable{
		{Name: "OPTIONAL_VAR", Required: false},
	})
	env := map[string]string{}

	report, err := validator.Validate(env, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasErrors() {
		t.Error("optional missing variable should not cause errors")
	}
}
