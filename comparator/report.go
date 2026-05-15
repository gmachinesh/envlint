package comparator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ReportText returns a human-readable summary of the comparison result.
func ReportText(r Result) string {
	if !r.HasDiff() {
		return "No differences found.\n"
	}

	var sb strings.Builder
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			fmt.Fprintf(&sb, "+ %s=%s\n", c.Key, c.NewValue)
		case Removed:
			fmt.Fprintf(&sb, "- %s=%s\n", c.Key, c.OldValue)
		case Modified:
			fmt.Fprintf(&sb, "~ %s: %q -> %q\n", c.Key, c.OldValue, c.NewValue)
		}
	}
	return sb.String()
}

// ReportJSON returns a JSON-encoded representation of the comparison result.
func ReportJSON(r Result) (string, error) {
	type jsonChange struct {
		Key      string `json:"key"`
		Kind     string `json:"kind"`
		OldValue string `json:"old_value,omitempty"`
		NewValue string `json:"new_value,omitempty"`
	}

	var out []jsonChange
	for _, c := range r.Changes {
		if c.Kind == Unchanged {
			continue
		}
		out = append(out, jsonChange{
			Key:      c.Key,
			Kind:     string(c.Kind),
			OldValue: c.OldValue,
			NewValue: c.NewValue,
		})
	}

	if out == nil {
		out = []jsonChange{}
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
