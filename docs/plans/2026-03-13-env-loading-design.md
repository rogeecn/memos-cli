# Environment Loading Design

## Background

The Go CLI currently reads configuration only from process environment variables via `internal/config/config.go`. The repository already documents `.env`-based configuration in `README.md`, but the Go implementation does not yet honor a `.env` file in the current working directory.

## Goal

Add `godotenv` support so the CLI can read configuration from the current directory `.env` file while preserving standard CLI behavior where explicit system environment variables take precedence.

## Chosen Approach

Use `godotenv.Read()` inside `config.LoadFromEnv()` to parse the current directory `.env` file into an in-memory map. Then resolve each config value with this precedence:

1. `os.Getenv(key)` if non-empty
2. value from `.env` if present
3. empty string

This keeps all config consumers on the same code path without mutating the process environment.

## Why This Approach

- Keeps configuration behavior centralized in one place.
- Avoids hidden global side effects from `godotenv.Load()`.
- Makes tests straightforward because `.env` parsing stays local to config loading.
- Matches expected CLI precedence: shell-provided environment wins over local defaults.

## Scope

- Modify `internal/config/config.go` to load `.env` from the current working directory.
- Add config tests for `.env` fallback and environment override behavior.
- Update `README.md` to explicitly describe the precedence rule.
- Add `github.com/joho/godotenv` to Go module dependencies.

## Error Handling

- Missing `.env` file should not be treated as an error.
- Parsing errors in an existing `.env` file should not crash all commands unexpectedly; the implementation should prefer safe behavior and keep validation focused on required config values.
- Validation rules for `MEMOS_URL`, `MEMOS_API_KEY`, and `MEMOS_ADMIN_API_KEY` remain unchanged.

## Testing Strategy

Use TDD in `internal/config/config_test.go`:

1. Write a failing test proving `.env` values are loaded when environment variables are absent.
2. Write a failing test proving process environment variables override `.env` values.
3. Run the targeted test command to confirm failure.
4. Implement the smallest production change to pass.
5. Re-run targeted tests, then `go test ./...`.

## Success Criteria

- Running the CLI in a directory with a `.env` file works without exporting variables first.
- Exported environment variables still override `.env` values.
- Existing commands that rely on `config.LoadFromEnv()` keep working without per-command changes.
