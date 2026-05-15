// Package resolver resolves the final effective value of an environment
// variable by consulting a priority-ordered chain of sources:
//
//  1. The caller-supplied env map (e.g. loaded from a .env file).
//  2. The host OS environment (os.Getenv), when FallbackToOS is enabled.
//  3. The schema-defined default value, when FallbackToDefault is enabled.
//
// Basic usage:
//
//	env, _ := dotenv.Load(".env")
//	s, _   := schema.Load("schema.yaml")
//
//	opts := resolver.DefaultOptions()
//	r    := resolver.Resolve("DATABASE_URL", env, s, opts)
//	if r.Found {
//		fmt.Println(r.Value, "from", r.Source)
//	}
//
// To resolve every key declared in the schema at once use ResolveAll.
// MustResolve returns an error when the key cannot be found in any source.
package resolver
