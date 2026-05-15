package scorer_test

import (
	"testing"

	"github.com/user/envlint/scorer"
	"github.com/user/envlint/validator"
)

func makeResults(severities ...string) []validator.Result {
	var out []validator.Result
	for _, s := range severities {
		out = append(out, validator.Result{Severity: s})
	}
	return out
}

func TestCompute_NoResults(t *testing.T) {
	s := scorer.Compute(nil, scorer.DefaultWeights)
	if s.Value != 100 {
		t.Fatalf("expected 100, got %d", s.Value)
	}
	if s.Penalties != 0 {
		t.Fatalf("expected 0 penalties, got %d", s.Penalties)
	}
}

func TestCompute_OnlyErrors(t *testing.T) {
	s := scorer.Compute(makeResults("error", "error"), scorer.DefaultWeights)
	// 2 errors * 10 = 20 penalty => score 80
	if s.Value != 80 {
		t.Fatalf("expected 80, got %d", s.Value)
	}
	if s.ErrorCount != 2 {
		t.Fatalf("expected ErrorCount 2, got %d", s.ErrorCount)
	}
}

func TestCompute_MixedSeverities(t *testing.T) {
	results := makeResults("error", "warning", "info")
	s := scorer.Compute(results, scorer.DefaultWeights)
	// 10 + 3 + 1 = 14 penalty => score 86
	if s.Value != 86 {
		t.Fatalf("expected 86, got %d", s.Value)
	}
	if s.WarningCount != 1 || s.InfoCount != 1 {
		t.Fatalf("unexpected counts: warn=%d info=%d", s.WarningCount, s.InfoCount)
	}
}

func TestCompute_ClampedToZero(t *testing.T) {
	// 15 errors * 10 = 150 penalty, clamped to 0
	var sev []string
	for i := 0; i < 15; i++ {
		sev = append(sev, "error")
	}
	s := scorer.Compute(makeResults(sev...), scorer.DefaultWeights)
	if s.Value != 0 {
		t.Fatalf("expected 0, got %d", s.Value)
	}
}

func TestCompute_CustomWeights(t *testing.T) {
	w := scorer.Weights{Error: 5, Warning: 1, Info: 0}
	s := scorer.Compute(makeResults("error"), w)
	if s.Value != 95 {
		t.Fatalf("expected 95, got %d", s.Value)
	}
}

func TestGrade(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{100, "A"}, {90, "A"}, {89, "B"}, {75, "B"},
		{74, "C"}, {60, "C"}, {59, "D"}, {40, "D"},
		{39, "F"}, {0, "F"},
	}
	for _, tc := range cases {
		got := scorer.Grade(tc.score)
		if got != tc.want {
			t.Errorf("Grade(%d) = %q, want %q", tc.score, got, tc.want)
		}
	}
}
