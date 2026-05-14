package validator

// Result holds the validation outcome for a single schema key.
type Result struct {
	// Key is the environment variable name.
	Key string

	// Error contains the validation failure message, or empty string on success.
	Error string
}

// HasError returns true when the result represents a validation failure.
func (r Result) HasError() bool {
	return r.Error != ""
}

// Summary returns a slice of Results that contain errors only.
func Summary(results []Result) []Result {
	var errs []Result
	for _, r := range results {
		if r.HasError() {
			errs = append(errs, r)
		}
	}
	return errs
}
