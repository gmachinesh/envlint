package masker_test

import (
	"fmt"

	"github.com/yourorg/envlint/masker"
)

func ExampleMask() {
	opts := masker.DefaultOptions()
	fmt.Println(masker.Mask("supersecret", opts))
	// Output: su*******et
}

func ExampleMaskKeys() {
	env := map[string]string{
		"API_KEY": "abc123xyz",
		"APP_ENV": "production",
	}
	opts := masker.DefaultOptions()
	masked := masker.MaskKeys(env, []string{"API_KEY"}, opts)
	fmt.Println(masked["APP_ENV"])
	// Output: production
}
