// Package masker provides utilities for partially masking sensitive
// environment variable values in output, preserving enough context
// for debugging while avoiding full secret exposure.
package masker

import "strings"

// Options controls how values are masked.
type Options struct {
	// ShowPrefix is the number of leading characters to reveal.
	ShowPrefix int
	// ShowSuffix is the number of trailing characters to reveal.
	ShowSuffix int
	// MaskChar is the character used for masking (default '*').
	MaskChar rune
	// MaskLen is the fixed length of the mask segment (0 = match original).
	MaskLen int
}

// DefaultOptions returns sensible masking defaults.
func DefaultOptions() Options {
	return Options{
		ShowPrefix: 2,
		ShowSuffix: 2,
		MaskChar:   '*',
		MaskLen:    0,
	}
}

// Mask partially obscures a single value according to opts.
// Values shorter than ShowPrefix+ShowSuffix are fully masked.
func Mask(value string, opts Options) string {
	if opts.MaskChar == 0 {
		opts.MaskChar = '*'
	}
	runes := []rune(value)
	total := len(runes)
	visible := opts.ShowPrefix + opts.ShowSuffix

	if total <= visible || total == 0 {
		maskLen := opts.MaskLen
		if maskLen == 0 {
			maskLen = total
			if maskLen == 0 {
				maskLen = 4
			}
		}
		return strings.Repeat(string(opts.MaskChar), maskLen)
	}

	prefix := string(runes[:opts.ShowPrefix])
	suffix := string(runes[total-opts.ShowSuffix:])

	maskLen := opts.MaskLen
	if maskLen == 0 {
		maskLen = total - visible
	}

	return prefix + strings.Repeat(string(opts.MaskChar), maskLen) + suffix
}

// MaskAll applies Mask to every value in the provided map, returning a new map.
func MaskAll(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = Mask(v, opts)
	}
	return out
}

// MaskKeys applies Mask only to the specified keys; all other keys are copied
// unchanged.
func MaskKeys(env map[string]string, keys []string, opts Options) map[string]string {
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		if _, sensitive := set[k]; sensitive {
			out[k] = Mask(v, opts)
		} else {
			out[k] = v
		}
	}
	return out
}
