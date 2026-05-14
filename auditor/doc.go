// Package auditor detects variables present in a .env file that have
// no corresponding definition in the schema. These "unknown" keys may
// indicate stale configuration, copy-paste errors, or undocumented
// secrets that should be reviewed.
//
// Usage:
//
//	results := auditor.Audit(envMap, schema)
//	for _, r := range results {
//		fmt.Printf("[%s] %s\n", r.Severity, r.Message)
//	}
//
// Findings carry a severity of "warning" by default since unknown keys
// are not necessarily errors — they may be application-specific overrides
// that the schema author has not yet documented.
package auditor
