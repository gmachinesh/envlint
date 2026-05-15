package coercer_test

import (
	"testing"

	"github.com/nicholasgasior/envlint/coercer"
)

func TestCoerce_TrimSpace(t *testing.T) {
	opts := coercer.DefaultOptions()
	got := coercer.Coerce("  hello  ", opts)
	if got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestCoerce_StripDoubleQuotes(t *testing.T) {
	opts := coercer.DefaultOptions()
	got := coercer.Coerce(`"my value"`, opts)
	if got != "my value" {
		t.Fatalf("expected 'my value', got %q", got)
	}
}

func TestCoerce_StripSingleQuotes(t *testing.T) {
	opts := coercer.DefaultOptions()
	got := coercer.Coerce("'my value'", opts)
	if got != "my value" {
		t.Fatalf("expected 'my value', got %q", got)
	}
}

func TestCoerce_NormaliseBool_Yes(t *testing.T) {
	opts := coercer.DefaultOptions()
	for _, input := range []string{"yes", "YES", "Yes", "1", "on", "ON"} {
		got := coercer.Coerce(input, opts)
		if got != "true" {
			t.Errorf("input %q: expected 'true', got %q", input, got)
		}
	}
}

func TestCoerce_NormaliseBool_No(t *testing.T) {
	opts := coercer.DefaultOptions()
	for _, input := range []string{"no", "NO", "0", "off", "OFF"} {
		got := coercer.Coerce(input, opts)
		if got != "false" {
			t.Errorf("input %q: expected 'false', got %q", input, got)
		}
	}
}

func TestCoerce_NoNormaliseBool(t *testing.T) {
	opts := coercer.DefaultOptions()
	opts.NormaliseBool = false
	got := coercer.Coerce("yes", opts)
	if got != "yes" {
		t.Fatalf("expected 'yes', got %q", got)
	}
}

func TestCoerce_NoStripQuotes(t *testing.T) {
	opts := coercer.DefaultOptions()
	opts.StripQuotes = false
	input := `"quoted"`
	got := coercer.Coerce(input, opts)
	if got != input {
		t.Fatalf("expected %q unchanged, got %q", input, got)
	}
}

func TestCoerceAll_AppliestoAllKeys(t *testing.T) {
	env := map[string]string{
		"A": "  hello  ",
		"B": "yes",
		"C": "'quoted'",
	}
	opts := coercer.DefaultOptions()
	out := coercer.CoerceAll(env, opts)

	if out["A"] != "hello" {
		t.Errorf("A: expected 'hello', got %q", out["A"])
	}
	if out["B"] != "true" {
		t.Errorf("B: expected 'true', got %q", out["B"])
	}
	if out["C"] != "quoted" {
		t.Errorf("C: expected 'quoted', got %q", out["C"])
	}
}

func TestCoerceAll_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"K": "  val  "}
	coercer.CoerceAll(env, coercer.DefaultOptions())
	if env["K"] != "  val  " {
		t.Fatal("original map was mutated")
	}
}
