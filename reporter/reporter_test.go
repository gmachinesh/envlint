package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlint/reporter"
	"github.com/user/envlint/validator"
)

func makeResults(pairs ...string) []validator.Result {
	var results []validator.Result
	for i := 0; i+1 < len(pairs); i += 2 {
		results = append(results, validator.Result{Key: pairs[i], Error: pairs[i+1]})
	}
	return results
}

func TestReportText_NoErrors(t *testing.T) {
	var buf bytes.Buffer
	r := &reporter.Reporter{Writer: &buf, Format: reporter.FormatText}
	results := makeResults("PORT", "", "HOST", "")
	hasErrors := r.Report(results)
	if hasErrors {
		t.Error("expected no errors")
	}
	out := buf.String()
	if !strings.Contains(out, "[OK]    PORT") {
		t.Errorf("expected OK for PORT, got: %s", out)
	}
}

func TestReportText_WithErrors(t *testing.T) {
	var buf bytes.Buffer
	r := &reporter.Reporter{Writer: &buf, Format: reporter.FormatText}
	results := makeResults("PORT", "missing required variable", "HOST", "")
	hasErrors := r.Report(results)
	if !hasErrors {
		t.Error("expected errors")
	}
	out := buf.String()
	if !strings.Contains(out, "[ERROR] PORT") {
		t.Errorf("expected ERROR for PORT, got: %s", out)
	}
	if !strings.Contains(out, "[OK]    HOST") {
		t.Errorf("expected OK for HOST, got: %s", out)
	}
}

func TestReportJSON_NoErrors(t *testing.T) {
	var buf bytes.Buffer
	r := &reporter.Reporter{Writer: &buf, Format: reporter.FormatJSON}
	results := makeResults("PORT", "")
	hasErrors := r.Report(results)
	if hasErrors {
		t.Error("expected no errors")
	}
	out := buf.String()
	if !strings.Contains(out, `"status":"ok"`) {
		t.Errorf("expected ok status in JSON, got: %s", out)
	}
}

func TestReportJSON_WithErrors(t *testing.T) {
	var buf bytes.Buffer
	r := &reporter.Reporter{Writer: &buf, Format: reporter.FormatJSON}
	results := makeResults("DB_URL", "missing required variable")
	hasErrors := r.Report(results)
	if !hasErrors {
		t.Error("expected errors")
	}
	out := buf.String()
	if !strings.Contains(out, `"status":"error"`) {
		t.Errorf("expected error status in JSON, got: %s", out)
	}
}

func TestReportText_Empty(t *testing.T) {
	var buf bytes.Buffer
	r := &reporter.Reporter{Writer: &buf, Format: reporter.FormatText}
	r.Report(nil)
	if !strings.Contains(buf.String(), "No variables") {
		t.Error("expected empty message")
	}
}
