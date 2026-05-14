# envlint

Validates `.env` files against a schema definition to catch missing or malformed variables in CI.

---

## Installation

```bash
go install github.com/yourname/envlint@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/yourname/envlint/releases).

---

## Usage

Define a schema file (`.env.schema`) describing the expected variables:

```ini
# .env.schema
DATABASE_URL=required,url
PORT=required,number
DEBUG=optional,bool
APP_NAME=required
```

Then validate your `.env` file against it:

```bash
envlint --schema .env.schema --env .env
```

Example output:

```
✗ DATABASE_URL: missing required variable
✗ PORT: expected number, got "abc"
✓ DEBUG: ok
✓ APP_NAME: ok

2 error(s) found.
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--schema` | `.env.schema` | Path to the schema file |
| `--env` | `.env` | Path to the `.env` file to validate |
| `--quiet` | `false` | Only print errors |

### CI Integration

```yaml
# .github/workflows/lint.yml
- name: Validate .env
  run: envlint --schema .env.schema --env .env.ci
```

---

## License

MIT © [yourname](https://github.com/yourname)