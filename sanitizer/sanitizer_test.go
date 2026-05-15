package sanitizer_test

import (
	"testing"

	"github.com/user/envlint/sanitizer"
)

func TestSanitize_DoesNotMutateOriginal(t *testing.T) {
	orig := map[string]string{"  KEY  ": "  value  "}
	opts := sanitizer.DefaultOptions()
	_ = sanitizer.Sanitize(orig, opts)
	if _, ok := orig["KEY"]; ok {
		t.Fatal("original map should not be mutated")
	}
}

func TestSanitize_TrimSpaceKey(t *testing.T) {
	env := map[string]string{"  FOO  ": "bar"}
	out := sanitizer.Sanitize(env, sanitizer.DefaultOptions())
	if _, ok := out["FOO"]; !ok {
		t.Fatalf("expected trimmed key 'FOO', got %v", out)
	}
}

func TestSanitize_TrimSpaceValue(t *testing.T) {
	env := map[string]string{"KEY": "  hello world  "}
	out := sanitizer.Sanitize(env, sanitizer.DefaultOptions())
	if got := out["KEY"]; got != "hello world" {
		t.Fatalf("expected 'hello world', got %q", got)
	}
}

func TestSanitize_RemoveControlChars(t *testing.T) {
	env := map[string]string{"KEY": "val\x01ue\x07"}
	out := sanitizer.Sanitize(env, sanitizer.DefaultOptions())
	if got := out["KEY"]; got != "value" {
		t.Fatalf("expected 'value', got %q", got)
	}
}

func TestSanitize_NormalizeNewlines(t *testing.T) {
	env := map[string]string{"KEY": "line1\r\nline2\rline3"}
	out := sanitizer.Sanitize(env, sanitizer.DefaultOptions())
	if got := out["KEY"]; got != "line1\nline2\nline3" {
		t.Fatalf("unexpected value: %q", got)
	}
}

func TestSanitize_CollapseInternalSpace(t *testing.T) {
	opts := sanitizer.DefaultOptions()
	opts.CollapseInternalSpace = true
	env := map[string]string{"MY  KEY": "value"}
	out := sanitizer.Sanitize(env, opts)
	if _, ok := out["MY_KEY"]; !ok {
		t.Fatalf("expected collapsed key 'MY_KEY', got %v", out)
	}
}

func TestSanitize_EmptyKeyDropped(t *testing.T) {
	opts := sanitizer.DefaultOptions()
	env := map[string]string{"   ": "value"}
	out := sanitizer.Sanitize(env, opts)
	if len(out) != 0 {
		t.Fatalf("expected empty result, got %v", out)
	}
}

func TestSanitize_NoOpWhenDisabled(t *testing.T) {
	opts := sanitizer.Options{}
	env := map[string]string{"  KEY  ": "  val  "}
	out := sanitizer.Sanitize(env, opts)
	if _, ok := out["  KEY  "]; !ok {
		t.Fatal("key should not be trimmed when TrimSpace is false")
	}
	if got := out["  KEY  "]; got != "  val  " {
		t.Fatalf("value should not be trimmed, got %q", got)
	}
}
