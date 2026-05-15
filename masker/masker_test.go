package masker_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlint/masker"
)

func TestMask_ShortValue_FullyMasked(t *testing.T) {
	opts := masker.DefaultOptions() // prefix=2, suffix=2
	result := masker.Mask("ab", opts)
	for _, ch := range result {
		if ch != '*' {
			t.Errorf("expected all mask chars, got %q", result)
		}
	}
}

func TestMask_NormalValue(t *testing.T) {
	opts := masker.DefaultOptions()
	result := masker.Mask("supersecret", opts)
	if !strings.HasPrefix(result, "su") {
		t.Errorf("expected prefix 'su', got %q", result)
	}
	if !strings.HasSuffix(result, "et") {
		t.Errorf("expected suffix 'et', got %q", result)
	}
	if !strings.Contains(result, "*") {
		t.Errorf("expected mask chars in middle, got %q", result)
	}
}

func TestMask_FixedMaskLen(t *testing.T) {
	opts := masker.DefaultOptions()
	opts.MaskLen = 4
	result := masker.Mask("mysecretvalue", opts)
	// prefix(2) + mask(4) + suffix(2) = 8 chars
	if len(result) != 8 {
		t.Errorf("expected length 8, got %d (%q)", len(result), result)
	}
}

func TestMask_EmptyValue(t *testing.T) {
	opts := masker.DefaultOptions()
	result := masker.Mask("", opts)
	if result == "" {
		t.Error("expected non-empty mask for empty value")
	}
}

func TestMask_CustomMaskChar(t *testing.T) {
	opts := masker.DefaultOptions()
	opts.MaskChar = '#'
	result := masker.Mask("hello_world", opts)
	if strings.Contains(result, "*") {
		t.Errorf("expected '#' mask char, got %q", result)
	}
	if !strings.Contains(result, "#") {
		t.Errorf("expected '#' in result, got %q", result)
	}
}

func TestMaskAll_MasksEveryValue(t *testing.T) {
	env := map[string]string{
		"KEY1": "password123",
		"KEY2": "apitoken456",
	}
	opts := masker.DefaultOptions()
	result := masker.MaskAll(env, opts)
	for k, v := range result {
		if v == env[k] {
			t.Errorf("key %s: value was not masked", k)
		}
	}
}

func TestMaskKeys_OnlyMasksSpecified(t *testing.T) {
	env := map[string]string{
		"SECRET": "topsecret",
		"PUBLIC": "openvalue",
	}
	opts := masker.DefaultOptions()
	result := masker.MaskKeys(env, []string{"SECRET"}, opts)
	if result["SECRET"] == env["SECRET"] {
		t.Error("SECRET should have been masked")
	}
	if result["PUBLIC"] != env["PUBLIC"] {
		t.Errorf("PUBLIC should be unchanged, got %q", result["PUBLIC"])
	}
}

func TestMaskKeys_EmptyKeyList(t *testing.T) {
	env := map[string]string{"A": "alpha", "B": "beta"}
	opts := masker.DefaultOptions()
	result := masker.MaskKeys(env, nil, opts)
	for k, v := range result {
		if v != env[k] {
			t.Errorf("key %s should be unchanged", k)
		}
	}
}
