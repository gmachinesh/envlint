// Package merger combines multiple .env files into a single map,
// with later files taking precedence over earlier ones.
package merger

import (
	"fmt"

	"github.com/user/envlint/dotenv"
)

// Options controls how merging behaves.
type Options struct {
	// FailOnMissing causes Merge to return an error if any file cannot be loaded.
	FailOnMissing bool
}

// DefaultOptions returns sensible defaults for merging.
func DefaultOptions() Options {
	return Options{
		FailOnMissing: true,
	}
}

// Merge loads each file in order and merges them into a single map.
// Keys from later files overwrite keys from earlier files.
func Merge(paths []string, opts Options) (map[string]string, error) {
	result := make(map[string]string)

	for _, p := range paths {
		env, err := dotenv.Load(p)
		if err != nil {
			if opts.FailOnMissing {
				return nil, fmt.Errorf("merger: loading %q: %w", p, err)
			}
			// skip missing/unreadable files silently
			continue
		}
		for k, v := range env {
			result[k] = v
		}
	}

	return result, nil
}

// Sources returns which file last defined each key, useful for diagnostics.
func Sources(paths []string, opts Options) (map[string]string, error) {
	sources := make(map[string]string)

	for _, p := range paths {
		env, err := dotenv.Load(p)
		if err != nil {
			if opts.FailOnMissing {
				return nil, fmt.Errorf("merger: loading %q: %w", p, err)
			}
			continue
		}
		for k := range env {
			sources[k] = p
		}
	}

	return sources, nil
}
