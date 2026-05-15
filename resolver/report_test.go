package resolver_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlint/resolver"
)

func sampleResults() []resolver.Result {
	return []resolver.Result{
		{Key: "APP_HOST", Value: "localhost", Source: "default", Found: true},
		{Key: "APP_PORT", Value: "8080", Source: "env", Found: true},
		{Key: "APP_SECRET", Source: "missing", Found: false},
	}
}

func TestReportText_NoResults(t *testing.T) {
	out := resolver.ReportText(nil)
	if !strings.Contains(out, "no variables resolved") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestReportText_WithResults(t *testing.T) {
	out := resolver.ReportText(sampleResults())
	if !strings.Contains(out, "2 resolved") {
		t.Errorf("expected resolved count, got: %q", out)
	}
	if !strings.Contains(out, "1 missing") {
		t.Errorf("expected missing count, got: %q", out)
	}
	if !strings.Contains(out, "APP_SECRET") {
		t.Errorf("expected APP_SECRET in output")
	}
}

func TestReportJSON_NoResults(t *testing.T) {
	out, err := resolver.ReportJSON(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "null" && out != "[]" {
		// json.MarshalIndent on empty slice returns []
		var arr []interface{}
		if jsonErr := json.Unmarshal([]byte(out), &arr); jsonErr != nil {
			t.Fatalf("invalid json: %v", jsonErr)
		}
	}
}

func TestReportJSON_WithResults(t *testing.T) {
	out, err := resolver.ReportJSON(sampleResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &arr); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(arr) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(arr))
	}
	found := false
	for _, entry := range arr {
		if entry["key"] == "APP_SECRET" && entry["found"] == false {
			found = true
		}
	}
	if !found {
		t.Error("expected APP_SECRET with found=false in JSON output")
	}
}
