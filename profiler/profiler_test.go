package profiler_test

import (
	"testing"

	"github.com/user/envlint/profiler"
)

func TestAnalyze_EmptyEnv(t *testing.T) {
	p := profiler.Analyze(map[string]string{})
	if p.Total != 0 {
		t.Fatalf("expected 0 total, got %d", p.Total)
	}
}

func TestAnalyze_BasicCounts(t *testing.T) {
	env := map[string]string{
		"APP_PORT":    "8080",
		"APP_DEBUG":   "true",
		"APP_URL":     "https://example.com",
		"APP_EMPTY":   "",
		"DB_HOST":     "localhost",
	}
	p := profiler.Analyze(env)

	if p.Total != 5 {
		t.Errorf("Total: want 5, got %d", p.Total)
	}
	if p.Empty != 1 {
		t.Errorf("Empty: want 1, got %d", p.Empty)
	}
	if p.Numeric != 1 {
		t.Errorf("Numeric: want 1, got %d", p.Numeric)
	}
	if p.Boolean != 1 {
		t.Errorf("Boolean: want 1, got %d", p.Boolean)
	}
	if p.URL != 1 {
		t.Errorf("URL: want 1, got %d", p.URL)
	}
}

func TestAnalyze_Prefixes(t *testing.T) {
	env := map[string]string{
		"APP_PORT":  "8080",
		"APP_DEBUG": "true",
		"DB_HOST":   "localhost",
		"NOPREFIX":  "value",
	}
	p := profiler.Analyze(env)

	if p.Prefixes["APP"] != 2 {
		t.Errorf("APP prefix count: want 2, got %d", p.Prefixes["APP"])
	}
	if p.Prefixes["DB"] != 1 {
		t.Errorf("DB prefix count: want 1, got %d", p.Prefixes["DB"])
	}
	if _, ok := p.Prefixes["NOPREFIX"]; ok {
		t.Error("NOPREFIX should not appear as a prefix")
	}
}

func TestAnalyze_LongestShortestKey(t *testing.T) {
	env := map[string]string{
		"A":             "1",
		"VERY_LONG_KEY": "2",
		"MID":           "3",
	}
	p := profiler.Analyze(env)

	if p.LongestKey != "VERY_LONG_KEY" {
		t.Errorf("LongestKey: want VERY_LONG_KEY, got %s", p.LongestKey)
	}
	if p.ShortestKey != "A" {
		t.Errorf("ShortestKey: want A, got %s", p.ShortestKey)
	}
}

func TestAnalyze_NegativeNumeric(t *testing.T) {
	env := map[string]string{"OFFSET": "-5"}
	p := profiler.Analyze(env)
	if p.Numeric != 1 {
		t.Errorf("expected negative number to be counted as numeric")
	}
}
