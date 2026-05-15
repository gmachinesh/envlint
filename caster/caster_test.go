package caster_test

import (
	"testing"
	"time"

	"github.com/envlint/envlint/caster"
)

func TestCastInt_Valid(t *testing.T) {
	r := caster.CastInt("PORT", "8080")
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if r.Value.(int) != 8080 {
		t.Fatalf("expected 8080, got %v", r.Value)
	}
}

func TestCastInt_Invalid(t *testing.T) {
	r := caster.CastInt("PORT", "abc")
	if r.Err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCastBool_True(t *testing.T) {
	for _, raw := range []string{"true", "1", "t", "TRUE"} {
		r := caster.CastBool("FLAG", raw)
		if r.Err != nil {
			t.Fatalf("unexpected error for %q: %v", raw, r.Err)
		}
		if !r.Value.(bool) {
			t.Fatalf("expected true for %q", raw)
		}
	}
}

func TestCastBool_Invalid(t *testing.T) {
	r := caster.CastBool("FLAG", "yes")
	if r.Err == nil {
		t.Fatal("expected error for 'yes'")
	}
}

func TestCastFloat_Valid(t *testing.T) {
	r := caster.CastFloat("RATIO", "3.14")
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	v := r.Value.(float64)
	if v < 3.13 || v > 3.15 {
		t.Fatalf("unexpected value: %v", v)
	}
}

func TestCastFloat_Invalid(t *testing.T) {
	r := caster.CastFloat("RATIO", "not-a-float")
	if r.Err == nil {
		t.Fatal("expected error")
	}
}

func TestCastDuration_Valid(t *testing.T) {
	r := caster.CastDuration("TIMEOUT", "30s")
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if r.Value.(time.Duration) != 30*time.Second {
		t.Fatalf("unexpected duration: %v", r.Value)
	}
}

func TestCastDuration_Invalid(t *testing.T) {
	r := caster.CastDuration("TIMEOUT", "30")
	if r.Err == nil {
		t.Fatal("expected error for bare number")
	}
}

func TestCast_String(t *testing.T) {
	r := caster.Cast("NAME", "hello", "string")
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if r.Value.(string) != "hello" {
		t.Fatalf("expected 'hello', got %v", r.Value)
	}
}

func TestCast_UnknownKind(t *testing.T) {
	r := caster.Cast("X", "v", "uuid")
	if r.Err == nil {
		t.Fatal("expected error for unknown kind")
	}
}
