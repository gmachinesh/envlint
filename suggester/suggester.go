// Package suggester provides hints for fixing common validation errors
// found in .env files, such as missing variables or pattern mismatches.
package suggester

import (
	"fmt"
	"strings"

	"github.com/user/envlint/validator"
)

// Suggestion holds a human-readable hint for a validation result.
type Suggestion struct {
	Key     string
	Message string
}

// Suggest returns a slice of Suggestion values derived from the given
// validation results. It inspects each result's severity and message to
// produce actionable advice.
func Suggest(results []validator.Result) []Suggestion {
	var suggestions []Suggestion

	for _, r := range results {
		if r.OK {
			continue
		}

		var hint string

		switch {
		case strings.Contains(r.Message, "missing"):
			hint = fmt.Sprintf("Add %q to your .env file. Check the schema for the expected format.", r.Key)

		case strings.Contains(r.Message, "pattern"):
			hint = fmt.Sprintf(
				"The value for %q does not match the required pattern. "+
					"Verify the format (e.g. URL, email, semver) defined in the schema.",
				r.Key,
			)

		case strings.Contains(r.Message, "empty"):
			hint = fmt.Sprintf("%q must not be empty. Provide a non-blank value.", r.Key)

		default:
			hint = fmt.Sprintf("Review the schema definition for %q and correct its value.", r.Key)
		}

		suggestions = append(suggestions, Suggestion{
			Key:     r.Key,
			Message: hint,
		})
	}

	return suggestions
}
