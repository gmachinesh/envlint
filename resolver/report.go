package resolver

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ReportText returns a human-readable summary of resolved results.
func ReportText(results []Result) string {
	if len(results) == 0 {
		return "resolver: no variables resolved\n"
	}
	var sb strings.Builder
	found, missing := 0, 0
	for _, r := range results {
		if r.Found {
			found++
			fmt.Fprintf(&sb, "  ✓ %-30s = %q (source: %s)\n", r.Key, r.Value, r.Source)
		} else {
			missing++
			fmt.Fprintf(&sb, "  ✗ %-30s   (missing)\n", r.Key)
		}
	}
	header := fmt.Sprintf("resolver: %d resolved, %d missing\n", found, missing)
	return header + sb.String()
}

// ReportJSON returns a JSON representation of resolved results.
func ReportJSON(results []Result) (string, error) {
	type jsonResult struct {
		Key    string `json:"key"`
		Value  string `json:"value,omitempty"`
		Source string `json:"source"`
		Found  bool   `json:"found"`
	}
	out := make([]jsonResult, len(results))
	for i, r := range results {
		out[i] = jsonResult{
			Key:    r.Key,
			Value:  r.Value,
			Source: r.Source,
			Found:  r.Found,
		}
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", fmt.Errorf("resolver: json marshal: %w", err)
	}
	return string(b), nil
}
