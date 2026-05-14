// Package auditor provides functionality to audit .env files for
// unused variables — keys present in the .env but not defined in the schema.
package auditor

import (
	"fmt"

	"github.com/user/envlint/schema"
)

// Result represents a single audit finding.
type Result struct {
	Key     string
	Message string
	Severity string
}

// Audit checks the provided env map against the schema and returns
// findings for keys that are present in env but absent from the schema.
func Audit(env map[string]string, s *schema.Schema) []Result {
	defined := make(map[string]bool, len(s.Vars))
	for _, v := range s.Vars {
		defined[v.Name] = true
	}

	var results []Result
	for key := range env {
		if !defined[key] {
			results = append(results, Result{
				Key:      key,
				Message:  fmt.Sprintf("key %q is not defined in schema", key),
				Severity: "warning",
			})
		}
	}
	return results
}

// HasUnknown returns true if any unknown keys were found.
func HasUnknown(results []Result) bool {
	return len(results) > 0
}
