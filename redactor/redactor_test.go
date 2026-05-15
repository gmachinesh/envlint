package redactor_test

import (
	"testing"

	"github.com/yourorg/envlint/redactor"
)

func TestRedact_NonSensitiveKeysUnchanged(t *testing.T) {
	env := map[string]string{
		"APP_ENV":  "production",
		"LOG_LEVEL": "info",
	}
	got := redactor.Redact(env, redactor.Options{})
	for k, want := range env {
		if got[k] != want {
			t.Errorf("key %s: got %q, want %q", k, got[k], want)
		}
	}
}

func TestRedact_SensitiveKeysRedacted(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD":  "supersecret",
		"API_KEY":      "abc123",
		"GITHUB_TOKEN": "ghp_xyz",
		"APP_SECRET":   "mysecret",
	}
	got := redactor.Redact(env, redactor.Options{})
	for k := range env {
		if got[k] != "[REDACTED]" {
			t.Errorf("key %s: expected [REDACTED], got %q", k, got[k])
		}
	}
}

func TestRedact_CustomPatterns(t *testing.T) {
	env := map[string]string{
		"STRIPE_KEY": "sk_live_123",
		"APP_NAME":   "envlint",
	}
	opts := redactor.Options{Patterns: []string{"STRIPE"}}
	got := redactor.Redact(env, opts)

	if got["STRIPE_KEY"] != "[REDACTED]" {
		t.Errorf("STRIPE_KEY should be redacted, got %q", got["STRIPE_KEY"])
	}
	if got["APP_NAME"] != "envlint" {
		t.Errorf("APP_NAME should not be redacted, got %q", got["APP_NAME"])
	}
}

func TestRedact_CaseInsensitiveMatch(t *testing.T) {
	env := map[string]string{
		"db_password": "secret",
		"Db_Password": "secret2",
	}
	got := redactor.Redact(env, redactor.Options{})
	for k := range env {
		if got[k] != "[REDACTED]" {
			t.Errorf("key %s: expected [REDACTED], got %q", k, got[k])
		}
	}
}

func TestRedact_EmptyEnv(t *testing.T) {
	got := redactor.Redact(map[string]string{}, redactor.Options{})
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"APP_ENV", false},
		{"AUTH_HEADER", true},
		{"LOG_LEVEL", false},
	}
	for _, tc := range cases {
		got := redactor.IsSensitive(tc.key, redactor.DefaultSensitivePatterns)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}
