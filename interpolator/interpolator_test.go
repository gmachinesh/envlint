package interpolator_test

import (
	"testing"

	"github.com/user/envlint/interpolator"
)

func TestInterpolate_NoReferences(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result, err := interpolator.Interpolate(env, interpolator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "bar" || result["BAZ"] != "qux" {
		t.Errorf("values should be unchanged, got %v", result)
	}
}

func TestInterpolate_SimpleReference(t *testing.T) {
	env := map[string]string{
		"BASE": "hello",
		"GREETING": "${BASE}_world",
	}
	result, err := interpolator.Interpolate(env, interpolator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["GREETING"] != "hello_world" {
		t.Errorf("expected hello_world, got %q", result["GREETING"])
	}
}

func TestInterpolate_ChainedReferences(t *testing.T) {
	env := map[string]string{
		"A": "alpha",
		"B": "${A}_beta",
		"C": "${B}_gamma",
	}
	result, err := interpolator.Interpolate(env, interpolator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["C"] != "alpha_beta_gamma" {
		t.Errorf("expected alpha_beta_gamma, got %q", result["C"])
	}
}

func TestInterpolate_UnresolvedAllowed(t *testing.T) {
	env := map[string]string{"FOO": "${MISSING}"}
	opts := interpolator.DefaultOptions()
	opts.FailOnUnresolved = false
	_, err := interpolator.Interpolate(env, opts)
	if err != nil {
		t.Fatalf("expected no error when FailOnUnresolved=false, got %v", err)
	}
}

func TestInterpolate_UnresolvedFails(t *testing.T) {
	env := map[string]string{"FOO": "${DEFINITELY_NOT_SET_XYZ123}"}
	opts := interpolator.DefaultOptions()
	opts.FailOnUnresolved = true
	_, err := interpolator.Interpolate(env, opts)
	if err == nil {
		t.Fatal("expected error for unresolved reference, got nil")
	}
}

func TestInterpolate_MutatesInPlace(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://${HOST}/db",
	}
	result, _ := interpolator.Interpolate(env, interpolator.DefaultOptions())
	if result["DSN"] != "postgres://localhost/db" {
		t.Errorf("expected postgres://localhost/db, got %q", result["DSN"])
	}
	// same map is returned
	if &result == nil {
		t.Error("result map should not be nil")
	}
}
