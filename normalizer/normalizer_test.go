package normalizer_test

import (
	"testing"

	"github.com/your-org/envlint/normalizer"
)

func TestNormalize_DefaultOptions(t *testing.T) {
	env := map[string]string{
		" db_host ": "  localhost  ",
		"api_key":   "Secret123",
	}
	opts := normalizer.DefaultOptions()
	got := normalizer.Normalize(env, opts)

	if v, ok := got["DB_HOST"]; !ok || v != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q (present=%v)", v, ok)
	}
	if v, ok := got["API_KEY"]; !ok || v != "Secret123" {
		t.Errorf("expected API_KEY=Secret123, got %q (present=%v)", v, ok)
	}
}

func TestNormalize_LowercaseValues(t *testing.T) {
	env := map[string]string{"MODE": "PRODUCTION"}
	opts := normalizer.Options{TrimSpace: true, UppercaseKeys: false, LowercaseValues: true}
	got := normalizer.Normalize(env, opts)
	if got["MODE"] != "production" {
		t.Errorf("expected production, got %q", got["MODE"])
	}
}

func TestNormalize_RemoveEmpty(t *testing.T) {
	env := map[string]string{
		"PRESENT": "value",
		"EMPTY":   "",
		"SPACES":  "   ",
	}
	opts := normalizer.Options{TrimSpace: true, RemoveEmpty: true}
	got := normalizer.Normalize(env, opts)

	if _, ok := got["EMPTY"]; ok {
		t.Error("expected EMPTY to be removed")
	}
	if _, ok := got["SPACES"]; ok {
		t.Error("expected SPACES to be removed after trim")
	}
	if got["PRESENT"] != "value" {
		t.Errorf("expected PRESENT=value, got %q", got["PRESENT"])
	}
}

func TestNormalize_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"key": "val"}
	opts := normalizer.DefaultOptions()
	normalizer.Normalize(env, opts)
	if _, ok := env["key"]; !ok {
		t.Error("original map was mutated")
	}
}

func TestNormalizeKey(t *testing.T) {
	cases := []struct{ in, want string }{
		{" db_host ", "DB_HOST"},
		{"api_key", "API_KEY"},
		{"ALREADY_UPPER", "ALREADY_UPPER"},
	}
	for _, tc := range cases {
		if got := normalizer.NormalizeKey(tc.in); got != tc.want {
			t.Errorf("NormalizeKey(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
