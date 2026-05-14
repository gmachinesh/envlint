package auditor

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ReportText writes a human-readable summary of audit results to w.
func ReportText(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "audit: no unknown keys found")
		return
	}
	fmt.Fprintf(w, "audit: %d unknown key(s) found\n", len(results))
	for _, r := range results {
		fmt.Fprintf(w, "  [%s] %s\n", strings.ToUpper(r.Severity), r.Message)
	}
}

// ReportJSON writes a JSON-encoded list of audit results to w.
func ReportJSON(w io.Writer, results []Result) error {
	type jsonResult struct {
		Key      string `json:"key"`
		Message  string `json:"message"`
		Severity string `json:"severity"`
	}
	out := make([]jsonResult, len(results))
	for i, r := range results {
		out[i] = jsonResult{Key: r.Key, Message: r.Message, Severity: r.Severity}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
