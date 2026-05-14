package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Variable describes a single expected environment variable.
type Variable struct {
	Name        string `yaml:"name"`
	Required    bool   `yaml:"required"`
	Pattern     string `yaml:"pattern"`
	Description string `yaml:"description"`
}

// Schema represents the full .env schema definition.
type Schema struct {
	Version   string     `yaml:"version"`
	Variables []Variable `yaml:"variables"`
}

// Load reads and parses a YAML schema file from the given path.
func Load(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading schema file: %w", err)
	}

	var s Schema
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing schema YAML: %w", err)
	}

	if len(s.Variables) == 0 {
		return nil, fmt.Errorf("schema contains no variables")
	}

	return &s, nil
}
