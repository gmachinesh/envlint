package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Env holds the parsed key-value pairs from a .env file.
type Env map[string]string

// Load reads a .env file and returns its key-value pairs.
// It skips blank lines and comments (lines starting with '#').
// It returns an error if the file cannot be opened or contains malformed lines.
func Load(path string) (Env, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("dotenv: cannot open file %q: %w", path, err)
	}
	defer f.Close()

	env := make(Env)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip optional inline comment
		if idx := strings.Index(line, " #"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("dotenv: malformed line %d in %q: %q", lineNum, path, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Strip surrounding quotes from value
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("dotenv: error reading %q: %w", path, err)
	}

	return env, nil
}
