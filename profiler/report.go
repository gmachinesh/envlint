package profiler

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ReportText writes a human-readable profile summary to w.
func ReportText(p Profile, w io.Writer) {
	fmt.Fprintf(w, "Profile Summary\n")
	fmt.Fprintf(w, "  Total variables : %d\n", p.Total)
	fmt.Fprintf(w, "  Empty values    : %d\n", p.Empty)
	fmt.Fprintf(w, "  Numeric values  : %d\n", p.Numeric)
	fmt.Fprintf(w, "  Boolean values  : %d\n", p.Boolean)
	fmt.Fprintf(w, "  URL values      : %d\n", p.URL)
	fmt.Fprintf(w, "  Longest key     : %s\n", p.LongestKey)
	fmt.Fprintf(w, "  Shortest key    : %s\n", p.ShortestKey)

	if len(p.Prefixes) == 0 {
		return
	}

	fmt.Fprintf(w, "  Prefixes:\n")
	prefixes := sortedPrefixes(p.Prefixes)
	for _, pr := range prefixes {
		fmt.Fprintf(w, "    %-20s %d\n", pr, p.Prefixes[pr])
	}
}

// ReportJSON encodes the profile as JSON into w.
func ReportJSON(p Profile, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

func sortedPrefixes(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
