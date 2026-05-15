package truncator_test

import (
	"strings"
	"testing"

	"github.com/user/envlint/truncator"
)

func TestTruncate_ShortString(t *testing.T) {
	opts := truncator.DefaultOptions()
	out := truncator.Truncate("hello", opts)
	if out != "hello" {
		t.Errorf("expected 'hello', got %q", out)
	}
}

func TestTruncate_ExactLength(t *testing.T) {
	opts := truncator.Options{MaxLen: 5, Suffix: "..."}
	out := truncator.Truncate("hello", opts)
	if out != "hello" {
		t.Errorf("expected 'hello', got %q", out)
	}
}

func TestTruncate_LongString(t *testing.T) {
	opts := truncator.Options{MaxLen: 5, Suffix: "..."}
	out := truncator.Truncate("hello world", opts)
	if out != "hello..." {
		t.Errorf("expected 'hello...', got %q", out)
	}
}

func TestTruncate_ZeroMaxLen(t *testing.T) {
	opts := truncator.Options{MaxLen: 0, Suffix: "..."}
	out := truncator.Truncate("hello world", opts)
	if out != "hello world" {
		t.Errorf("expected unchanged string, got %q", out)
	}
}

func TestTruncate_MultiByte(t *testing.T) {
	opts := truncator.Options{MaxLen: 3, Suffix: ""}
	// each emoji is multiple bytes but one rune
	out := truncator.Truncate("🔑🔒🔓🔐", opts)
	if out != "🔑🔒🔓" {
		t.Errorf("expected 3 emoji runes, got %q", out)
	}
}

func TestTruncateAll_OnlyValues(t *testing.T) {
	env := map[string]string{
		"SHORT_KEY": "tiny",
		"LONG_KEY":  strings.Repeat("x", 80),
	}
	opts := truncator.Options{MaxLen: 10, Suffix: "...", OnlyValues: true}
	out := truncator.TruncateAll(env, opts)

	if out["SHORT_KEY"] != "tiny" {
		t.Errorf("short value should be unchanged, got %q", out["SHORT_KEY"])
	}
	if len([]rune(out["LONG_KEY"])) != 13 { // 10 runes + 3 suffix chars
		t.Errorf("long value not truncated correctly: %q", out["LONG_KEY"])
	}
	if _, ok := out["LONG_KEY"]; !ok {
		t.Error("key should be preserved when OnlyValues=true")
	}
}

func TestTruncateAll_DoesNotMutateOriginal(t *testing.T) {
	original := map[string]string{"KEY": strings.Repeat("v", 100)}
	opts := truncator.DefaultOptions()
	_ = truncator.TruncateAll(original, opts)
	if len(original["KEY"]) != 100 {
		t.Error("original map was mutated")
	}
}

func TestTruncateAll_KeysTruncated(t *testing.T) {
	env := map[string]string{
		strings.Repeat("K", 20): "value",
	}
	opts := truncator.Options{MaxLen: 5, Suffix: "~", OnlyValues: false}
	out := truncator.TruncateAll(env, opts)
	for k := range out {
		if k != "KKKKK~" {
			t.Errorf("expected truncated key 'KKKKK~', got %q", k)
		}
	}
}
