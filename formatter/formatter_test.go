package formatter_test

import (
	"testing"

	"github.com/user/envlint/formatter"
)

func TestFormatKey_Default(t *testing.T) {
	opts := formatter.DefaultOptions()
	result := formatter.FormatKey("MY_VAR", opts)
	if result != "MY_VAR" {
		t.Errorf("expected MY_VAR, got %s", result)
	}
}

func TestFormatKey_Uppercase(t *testing.T) {
	opts := formatter.Options{Style: formatter.StyleUppercase}
	result := formatter.FormatKey("my_var", opts)
	if result != "MY_VAR" {
		t.Errorf("expected MY_VAR, got %s", result)
	}
}

func TestFormatKey_Lowercase(t *testing.T) {
	opts := formatter.Options{Style: formatter.StyleLowercase}
	result := formatter.FormatKey("MY_VAR", opts)
	if result != "my_var" {
		t.Errorf("expected my_var, got %s", result)
	}
}

func TestFormatKey_SnakeCase(t *testing.T) {
	opts := formatter.Options{Style: formatter.StyleSnakeCase}
	result := formatter.FormatKey("My-Var Name", opts)
	if result != "my_var_name" {
		t.Errorf("expected my_var_name, got %s", result)
	}
}

func TestFormatKey_WithPrefix(t *testing.T) {
	opts := formatter.Options{Prefix: "APP_"}
	result := formatter.FormatKey("PORT", opts)
	if result != "APP_PORT" {
		t.Errorf("expected APP_PORT, got %s", result)
	}
}

func TestFormatValue_Plain(t *testing.T) {
	opts := formatter.DefaultOptions()
	result := formatter.FormatValue("hello", opts)
	if result != "hello" {
		t.Errorf("expected hello, got %s", result)
	}
}

func TestFormatValue_Masked(t *testing.T) {
	opts := formatter.Options{MaskValues: true}
	result := formatter.FormatValue("secret123", opts)
	if result != "***" {
		t.Errorf("expected ***, got %s", result)
	}
}

func TestFormatValue_MaskedEmpty(t *testing.T) {
	opts := formatter.Options{MaskValues: true}
	result := formatter.FormatValue("", opts)
	if result != "" {
		t.Errorf("expected empty string, got %s", result)
	}
}

func TestFormatValue_Truncated(t *testing.T) {
	opts := formatter.Options{MaxValueLen: 5}
	result := formatter.FormatValue("toolongvalue", opts)
	if result != "toolo..." {
		t.Errorf("expected 'toolo...', got %s", result)
	}
}

func TestFormatValue_NoTruncation(t *testing.T) {
	opts := formatter.Options{MaxValueLen: 20}
	result := formatter.FormatValue("short", opts)
	if result != "short" {
		t.Errorf("expected short, got %s", result)
	}
}
