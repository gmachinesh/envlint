package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// VarType represents the expected type of an environment variable.
type VarType string

const (
	TypeString VarType = "string"
	TypeInt    VarType = "int"
	TypeBool   VarType = "bool"
	TypeURL    VarType = "url"
)

// VarSchema defines the rules for a single environment variable.
type VarSchema struct {
	Required    bool    `yaml:"required"`
	Type        VarType `yaml:"type"`
	Pattern     string  `yaml:"pattern"`
	Description string  `yaml:"description"`
}

// Schema represents the full schema definition loaded from a YAML file.
type Schema struct {
	Vars map[string]VarSchema `yaml:"vars"`
}

// Load reads and parses a schema YAML file from the given path.
func Load(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading schema file %q: %w", path, err)
	}

	var s Schema
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing schema file %q: %w", path, err)
	}

	if s.Vars == nil {
		s.Vars = make(map[string]VarSchema)
	}

	return &s, nil
}
