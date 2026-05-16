// Package digester computes deterministic cryptographic digests over env maps.
//
// It is useful for:
//   - Detecting whether an .env file has changed between CI runs.
//   - Generating a cache key derived from the current environment.
//   - Verifying integrity of an env snapshot at rest.
//
// # Algorithms
//
// Two algorithms are supported:
//
//	AlgoSHA256  – standard SHA-256 (recommended for integrity checks)
//	AlgoFNV     – fast non-cryptographic FNV-1a 64-bit (suitable for cache keys)
//
// # Usage
//
//	env := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
//
//	result, err := digester.Digest(env, digester.AlgoSHA256)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.Hex) // e.g. "a3f1..."
//
//	// Quick equality check
//	if !digester.Equal(envA, envB) {
//	    fmt.Println("environments differ")
//	}
package digester
