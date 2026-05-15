// Package comparator provides utilities for comparing two env var maps and
// producing structured Change records that describe what was added, removed,
// or modified between them.
//
// Basic usage:
//
//	base := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
//	head := map[string]string{"DB_HOST": "prod.db", "PORT": "5432", "DEBUG": "true"}
//
//	result := comparator.Compare(base, head)
//	if result.HasDiff() {
//		fmt.Print(comparator.ReportText(result))
//	}
//
// Output formats:
//
//	ReportText – human-readable unified-diff-style output.
//	ReportJSON – machine-readable JSON array of changes.
package comparator
