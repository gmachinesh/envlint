// Package sorter provides deterministic ordering for environment variable
// keys and validation results produced by envlint.
//
// Three ordering modes are available:
//
//   - OrderAsc  – alphabetical ascending (A → Z), case-insensitive.
//   - OrderDesc – alphabetical descending (Z → A), case-insensitive.
//   - OrderSeverity – groups results by severity (error → warning → info → ok)
//     with a secondary alphabetical sort by key within each group.
//
// Sorting functions never mutate the input slice; they always return a new
// copy so the caller retains ownership of the original data.
//
// Example:
//
//	import "github.com/nicholasgasior/envlint/sorter"
//
//	keys   := []string{"PORT", "APP_ENV", "DB_URL"}
//	sorted := sorter.SortKeys(keys, sorter.DefaultOptions())
//	// sorted → ["APP_ENV", "DB_URL", "PORT"]
package sorter
