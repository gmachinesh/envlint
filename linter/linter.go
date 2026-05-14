package linter

import (
	"fmt"

	"github.com/user/envlint/dotenv"
	"github.com/user/envlint/expander"
	"github.com/user/envlint/filter"
	"github.com/user/envlint/schema"
	"github.com/user/envlint/validator"
)

// Options controls linter behaviour.
type Options struct {
	SchemaPath   string
	EnvPath      string
	Severity     string
	Prefix       string
	ExpandVars   bool
}

// Result holds the combined output of a lint run.
type Result struct {
	Results []validator.Result
	Summary validator.Summary
}

// Run loads the schema and env file, optionally expands variables,
// validates, and applies filters before returning a Result.
func Run(opts Options) (*Result, error) {
	sc, err := schema.Load(opts.SchemaPath)
	if err != nil {
		return nil, fmt.Errorf("loading schema: %w", err)
	}

	env, err := dotenv.Load(opts.EnvPath)
	if err != nil {
		return nil, fmt.Errorf("loading env file: %w", err)
	}

	if opts.ExpandVars {
		expanded := make(map[string]string, len(env))
		for k, v := range env {
			expanded[k] = expander.Expand(v, env)
		}
		env = expanded
	}

	results := validator.Validate(sc, env)

	fo := filter.Options{
		Severity: opts.Severity,
		Prefix:   opts.Prefix,
	}
	filtered := filter.Apply(results, fo)

	sum := validator.Summary{}
	for _, r := range filtered {
		switch r.Severity {
		case "error":
			sum.Errors++
		case "warning":
			sum.Warnings++
		case "info":
			sum.Infos++
		}
	}
	sum.Total = len(filtered)

	return &Result{
		Results: filtered,
		Summary: sum,
	}, nil
}
