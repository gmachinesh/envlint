// Package dotenv parses .env files into a key/value map.
// Lines beginning with # are treated as comments.
// Inline comments (preceded by whitespace and #) are stripped.
// Quoted values (single or double) have their quotes removed.
package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/user/envlint/expander"
)

// Load reads the .env file at path and returns a map of key/value pairs.
// Variable references in values are expanded using the parsed map and
// the OS environment.
func Load(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("dotenv: open %q: %w", path, err)
	}
	defer f.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip export prefix
		line = strings.TrimPrefix(line, "export ")

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("dotenv: line %d: malformed assignment %q", lineNum, line)
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		// Strip inline comment (only when value is unquoted)
		if len(val) == 0 || (val[0] != '"' && val[0] != '\'') {
			if ci := strings.Index(val, " #"); ci >= 0 {
				val = strings.TrimSpace(val[:ci])
			}
		}

		// Strip surrounding quotes
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		env[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("dotenv: scan %q: %w", path, err)
	}

	if expander.HasReference(fmt.Sprintf("%v", env)) {
		env = expander.Expand(env)
	}

	return env, nil
}
