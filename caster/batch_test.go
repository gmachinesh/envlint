package caster_test

import (
	"testing"

	"github.com/envlint/envlint/caster"
)

func TestCastAll_AllValid(t *testing.T) {
	inputs := []caster.BatchInput{
		{Key: "PORT", Raw: "3000", Kind: "int"},
		{Key: "DEBUG", Raw: "true", Kind: "bool"},
		{Key: "NAME", Raw: "app", Kind: "string"},
	}
	br := caster.CastAll(inputs)
	if br.HasErrors() {
		t.Fatalf("expected no errors, got %d failures", len(br.Failures))
	}
	if len(br.Successes) != 3 {
		t.Fatalf("expected 3 successes, got %d", len(br.Successes))
	}
}

func TestCastAll_PartialFailure(t *testing.T) {
	inputs := []caster.BatchInput{
		{Key: "PORT", Raw: "3000", Kind: "int"},
		{Key: "TIMEOUT", Raw: "notaduration", Kind: "duration"},
		{Key: "RATIO", Raw: "bad", Kind: "float"},
	}
	br := caster.CastAll(inputs)
	if !br.HasErrors() {
		t.Fatal("expected errors")
	}
	if len(br.Failures) != 2 {
		t.Fatalf("expected 2 failures, got %d", len(br.Failures))
	}
	if len(br.Successes) != 1 {
		t.Fatalf("expected 1 success, got %d", len(br.Successes))
	}
}

func TestCastAll_Empty(t *testing.T) {
	br := caster.CastAll(nil)
	if br.HasErrors() {
		t.Fatal("expected no errors on empty input")
	}
	if len(br.Successes) != 0 {
		t.Fatal("expected zero successes")
	}
}

func TestCastAll_AllFail(t *testing.T) {
	inputs := []caster.BatchInput{
		{Key: "A", Raw: "nope", Kind: "int"},
		{Key: "B", Raw: "nope", Kind: "float"},
	}
	br := caster.CastAll(inputs)
	if len(br.Failures) != 2 {
		t.Fatalf("expected 2 failures, got %d", len(br.Failures))
	}
	if len(br.Successes) != 0 {
		t.Fatalf("expected 0 successes, got %d", len(br.Successes))
	}
}
