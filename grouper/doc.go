// Package grouper partitions a slice of validator.Result values into named
// groups based on a chosen dimension.
//
// Three grouping strategies are supported:
//
//   - ByPrefix  — groups keys by the segment before the first underscore
//     (e.g. DB_HOST and DB_PASS → "db").
//   - BySeverity — groups by the result's Severity field ("error",
//     "warning", "info").
//   - ByStatus   — separates passing results ("ok") from failing ones
//     ("fail").
//
// All returned slices are in deterministic, alphabetically sorted order so
// that reports are reproducible across runs.
//
// Example:
//
//	groups := grouper.Apply(results, grouper.ByPrefix)
//	for _, g := range groups {
//		fmt.Printf("[%s] %d result(s)\n", g.Name, len(g.Results))
//	}
package grouper
