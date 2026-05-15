package profiler_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlint/profiler"
)

func sampleProfile() profiler.Profile {
	return profiler.Analyze(map[string]string{
		"APP_PORT":  "8080",
		"APP_DEBUG": "true",
		"DB_HOST":   "localhost",
		"DB_PASS":   "",
	})
}

func TestReportText_ContainsTotal(t *testing.T) {
	var buf bytes.Buffer
	profiler.ReportText(sampleProfile(), &buf)
	out := buf.String()

	if !strings.Contains(out, "Total variables") {
		t.Error("expected 'Total variables' in text report")
	}
	if !strings.Contains(out, "4") {
		t.Error("expected count 4 in text report")
	}
}

func TestReportText_ContainsPrefixes(t *testing.T) {
	var buf bytes.Buffer
	profiler.ReportText(sampleProfile(), &buf)
	out := buf.String()

	if !strings.Contains(out, "APP") {
		t.Error("expected prefix APP in text report")
	}
	if !strings.Contains(out, "DB") {
		t.Error("expected prefix DB in text report")
	}
}

func TestReportJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := profiler.ReportJSON(sampleProfile(), &buf); err != nil {
		t.Fatalf("ReportJSON error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if _, ok := out["Total"]; !ok {
		t.Error("expected 'Total' field in JSON output")
	}
}

func TestReportJSON_EmptyProfile(t *testing.T) {
	var buf bytes.Buffer
	p := profiler.Analyze(map[string]string{})
	if err := profiler.ReportJSON(p, &buf); err != nil {
		t.Fatalf("ReportJSON error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty JSON for empty profile")
	}
}
