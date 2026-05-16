package digester_test

import (
	"strings"
	"testing"

	"github.com/nicholasgasior/envlint/digester"
)

func TestDigest_SHA256_Deterministic(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r1, err := digester.Digest(env, digester.AlgoSHA256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r2, err := digester.Digest(env, digester.AlgoSHA256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r1.Hex != r2.Hex {
		t.Errorf("expected same digest, got %q and %q", r1.Hex, r2.Hex)
	}
	if r1.KeyCount != 2 {
		t.Errorf("expected KeyCount 2, got %d", r1.KeyCount)
	}
}

func TestDigest_SHA256_OrderIndependent(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2", "C": "3"}
	b := map[string]string{"C": "3", "A": "1", "B": "2"}
	ra, _ := digester.Digest(a, digester.AlgoSHA256)
	rb, _ := digester.Digest(b, digester.AlgoSHA256)
	if ra.Hex != rb.Hex {
		t.Errorf("order should not affect digest: %q vs %q", ra.Hex, rb.Hex)
	}
}

func TestDigest_SHA256_DifferentValues(t *testing.T) {
	a := map[string]string{"KEY": "value1"}
	b := map[string]string{"KEY": "value2"}
	ra, _ := digester.Digest(a, digester.AlgoSHA256)
	rb, _ := digester.Digest(b, digester.AlgoSHA256)
	if ra.Hex == rb.Hex {
		t.Error("expected different digests for different values")
	}
}

func TestDigest_FNV(t *testing.T) {
	env := map[string]string{"X": "y"}
	r, err := digester.Digest(env, digester.AlgoFNV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Hex) != 16 {
		t.Errorf("expected 16-char hex, got %q", r.Hex)
	}
	if r.Algorithm != digester.AlgoFNV {
		t.Errorf("expected AlgoFNV, got %q", r.Algorithm)
	}
}

func TestDigest_UnknownAlgorithm(t *testing.T) {
	_, err := digester.Digest(map[string]string{}, "md5")
	if err == nil {
		t.Error("expected error for unknown algorithm")
	}
	if !strings.Contains(err.Error(), "unknown algorithm") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEqual_SameMaps(t *testing.T) {
	env := map[string]string{"A": "1"}
	if !digester.Equal(env, env) {
		t.Error("expected Equal to return true for same map")
	}
}

func TestEqual_DifferentMaps(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"A": "2"}
	if digester.Equal(a, b) {
		t.Error("expected Equal to return false for different maps")
	}
}

func TestDigest_EmptyEnv(t *testing.T) {
	r, err := digester.Digest(map[string]string{}, digester.AlgoSHA256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.KeyCount != 0 {
		t.Errorf("expected KeyCount 0, got %d", r.KeyCount)
	}
	if r.Hex == "" {
		t.Error("expected non-empty hex for empty env")
	}
}
