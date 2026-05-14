package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envlint/schema"
)

// Result holds the outcome of a single variable validation.
type Result struct {
	Key     string
	Passed  bool
	Message string
}

// Report aggregates all validation results.
type Report struct {
	Results []Result
}

// HasErrors returns true if any result failed.
func (r *Report) HasErrors() bool {
	for _, res := range r.Results {
		if !res.Passed {
			return true
		}
	}
	return false
}

// Validate checks the provided env map against the schema.
func Validate(env map[string]string, s *schema.Schema) (*Report, error) {
	if s == nil {
		return nil, fmt.Errorf("schema must not be nil")
	}

	report := &Report{}

	for _, variable := range s.Variables {
		val, exists := env[variable.Name]

		if !exists || strings.TrimSpace(val) == "" {
			if variable.Required {
				report.Results = append(report.Results, Result{
					Key:     variable.Name,
					Passed:  false,
					Message: fmt.Sprintf("%s is required but missing or empty", variable.Name),
				})
			}
			continue
		}

		if variable.Pattern != "" {
			matched, err := regexp.MatchString(variable.Pattern, val)
			if err != nil {
				return nil, fmt.Errorf("invalid pattern for %s: %w", variable.Name, err)
			}
			if !matched {
				report.Results = append(report.Results, Result{
					Key:     variable.Name,
					Passed:  false,
					Message: fmt.Sprintf("%s value %q does not match pattern %q", variable.Name, val, variable.Pattern),
				})
				continue
			}
		}

		report.Results = append(report.Results, Result{
			Key:    variable.Name,
			Passed: true,
		})
	}

	return report, nil
}
