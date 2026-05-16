// Package digester computes stable hashes of env maps for change detection,
// caching, and integrity verification.
package digester

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

// Algorithm selects the hash representation.
type Algorithm string

const (
	AlgoSHA256 Algorithm = "sha256"
	AlgoFNV    Algorithm = "fnv"
)

// Result holds the computed digest for an env map.
type Result struct {
	Algorithm Algorithm
	Hex       string
	KeyCount  int
}

// Digest computes a deterministic hash over the provided env map.
// Keys are sorted before hashing so insertion order does not matter.
func Digest(env map[string]string, algo Algorithm) (Result, error) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	switch algo {
	case AlgoSHA256, "":
		return sha256Digest(env, keys)
	case AlgoFNV:
		return fnvDigest(env, keys)
	default:
		return Result{}, fmt.Errorf("digester: unknown algorithm %q", algo)
	}
}

// Equal returns true when two env maps produce the same digest.
func Equal(a, b map[string]string) bool {
	ra, err := Digest(a, AlgoSHA256)
	if err != nil {
		return false
	}
	rb, err := Digest(b, AlgoSHA256)
	if err != nil {
		return false
	}
	return ra.Hex == rb.Hex
}

func sha256Digest(env map[string]string, keys []string) (Result, error) {
	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, env[k])
	}
	return Result{
		Algorithm: AlgoSHA256,
		Hex:       hex.EncodeToString(h.Sum(nil)),
		KeyCount:  len(keys),
	}, nil
}

func fnvDigest(env map[string]string, keys []string) (Result, error) {
	var h uint64 = 14695981039346656037
	const prime uint64 = 1099511628211
	for _, k := range keys {
		for _, c := range []byte(k + "=" + env[k] + "\n") {
			h ^= uint64(c)
			h *= prime
		}
	}
	return Result{
		Algorithm: AlgoFNV,
		Hex:       fmt.Sprintf("%016x", h),
		KeyCount:  len(keys),
	}, nil
}
