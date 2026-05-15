package comparator_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/nicholasgasior/envlint/comparator"
)

func sampleResult() comparator.Result {
	return comparator.Compare(
		map[string]string{"OLD": "x", "SAME": "v"},
		map[string]string{"NEW": "y", "SAME": "v"},
	)
}

func TestReportText_NoDiff(t *testing.T) {
	r := comparator.Compare(
		map[string]string{"A": "1"},
		map[string]string{"A": "1"},
	)
	out := comparator.ReportText(r)
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected no-diff message, got: %s", out)
	}
}

func TestReportText_WithDiff(t *testing.T) {
	out := comparator.ReportText(sampleResult())
	if !strings.Contains(out, "+ NEW") {
		t.Errorf("expected added line, got: %s", out)
	}
	if !strings.Contains(out, "- OLD") {
		t.Errorf("expected removed line, got: %s", out)
	}
}

func TestReportJSON_NoDiff(t *testing.T) {
	r := comparator.Compare(
		map[string]string{"A": "1"},
		map[string]string{"A": "1"},
	)
	out, err := comparator.ReportJSON(r)
	if err != nil {
		t.Fatal(err)
	}
	var arr []interface{}
	if err := json.Unmarshal([]byte(out), &arr); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(arr) != 0 {
		t.Errorf("expected empty array, got %d items", len(arr))
	}
}

func TestReportJSON_WithDiff(t *testing.T) {
	out, err := comparator.ReportJSON(sampleResult())
	if err != nil {
		t.Fatal(err)
	}
	var arr []map[string]string
	if err := json.Unmarshal([]byte(out), &arr); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(arr) != 2 {
		t.Errorf("expected 2 changes, got %d", len(arr))
	}
}
