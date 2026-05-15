// Package masker provides partial-masking utilities for environment variable
// values.
//
// When displaying or logging environment variables it is often necessary to
// show enough of a value to confirm it is set correctly without exposing the
// full secret.  masker solves this by revealing a configurable number of
// leading and trailing characters while replacing the middle portion with a
// repeated mask character.
//
// # Basic usage
//
//	opts := masker.DefaultOptions()
//	masked := masker.Mask("supersecrettoken", opts)
//	// => "su***********en"
//
// # Masking a whole map
//
//	masked := masker.MaskAll(env, opts)
//
// # Masking only specific keys
//
//	masked := masker.MaskKeys(env, []string{"API_KEY", "DB_PASSWORD"}, opts)
package masker
