package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envlint/dotenv"
	"github.com/user/envlint/reporter"
	"github.com/user/envlint/schema"
	"github.com/user/envlint/validator"
)

func main() {
	envFile := flag.String("env", ".env", "path to the .env file")
	schemaFile := flag.String("schema", ".env.schema.yaml", "path to the schema YAML file")
	format := flag.String("format", "text", "output format: text or json")
	flag.Parse()

	sch, err := schema.Load(*schemaFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading schema: %v\n", err)
		os.Exit(2)
	}

	env, err := dotenv.Load(*envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading env file: %v\n", err)
		os.Exit(2)
	}

	results := validator.Validate(sch, env)

	fmt := reporter.FormatText
	if *format == "json" {
		fmt = reporter.FormatJSON
	}

	r := reporter.New(fmt)
	hasErrors := r.Report(results)
	if hasErrors {
		os.Exit(1)
	}
}
