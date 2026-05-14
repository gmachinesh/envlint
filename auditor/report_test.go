package auditor_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlint/auditor"
)

func sampleResults() []auditor.Result {
	return []auditor.Result{
		{Key: "OLD_SECRET", Message: `key "OLD_SECRET" is not defined in schema`, Severity: "warning"},
	}
}

func TestReportText_NoResults(t *testing.T) {
	var buf bytes.Buffer
	auditor.ReportText(&buf, nil)
	if !strings.Contains(buf.String(), "no unknown keys") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestReportText_WithResults(t *testing.T) {
	var buf bytes.Buffer
	auditor.ReportText(&buf, sampleResults())
	out := buf.String()
	if !strings.Contains(out, "1 unknown key") {
		t.Errorf("expected count in output, got: %q", out)
	}
	if !strings.Contains(out, "OLD_SECRET") {
		t.Errorf("expected key name in output, got: %q", out)
	}
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected severity in output, got: %q", out)
	}
}

func TestReportJSON_NoResults(t *testing.T) {
	var buf bytes.Buffer
	if err := auditor.ReportJSON(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty array, got %d items", len(out))
	}
}

func TestReportJSON_WithResults(t *testing.T) {
	var buf bytes.Buffer
	if err := auditor.ReportJSON(&buf, sampleResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 item, got %d", len(out))
	}
	if out[0]["key"] != "OLD_SECRET" {
		t.Errorf("expected key OLD_SECRET, got %v", out[0]["key"])
	}
}
