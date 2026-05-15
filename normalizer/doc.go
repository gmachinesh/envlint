// Package normalizer provides key and value normalisation for environment
// variable maps.
//
// Normalisation is applied before validation or comparison to ensure that
// minor formatting differences — such as inconsistent casing, surrounding
// whitespace, or empty values — do not cause spurious lint failures.
//
// Basic usage:
//
//	env := map[string]string{
//		" db_host ": "  localhost  ",
//	}
//
//	normalised := normalizer.Normalize(env, normalizer.DefaultOptions())
//	// normalised == map[string]string{"DB_HOST": "localhost"}
//
// Options can be customised to suit the project's conventions:
//
//	opts := normalizer.Options{
//		TrimSpace:     true,
//		UppercaseKeys: true,
//		RemoveEmpty:   true,
//	}
//	normalised := normalizer.Normalize(env, opts)
package normalizer
