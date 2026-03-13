# Godotenv Support Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add `godotenv` support so the Go CLI reads configuration from the current directory `.env` file while letting explicit system environment variables override `.env` values.

**Architecture:** Keep configuration loading centralized in `internal/config/config.go`. Parse `.env` into an in-memory map with `godotenv.Read()` and resolve each config field through a helper that prefers `os.Getenv()` over the parsed file value. This avoids mutating global process environment and keeps all commands on the same configuration path.

**Tech Stack:** Go 1.25, Cobra CLI, `github.com/joho/godotenv`, Go testing package.

---

### Task 1: Add failing config tests

**Files:**
- Modify: `internal/config/config_test.go`
- Test: `internal/config/config_test.go`

**Step 1: Write the failing test**

Add a test that creates a temporary working directory, writes a `.env` file containing:

```env
MEMOS_URL=http://from-dotenv
MEMOS_API_KEY=dotenv-key
MEMOS_ADMIN_API_KEY=dotenv-admin
DEFAULT_TAG=dotenv-tag
```

Then change into that directory, clear related environment variables, call `LoadFromEnv()`, and assert:

```go
if cfg.BaseURL != "http://from-dotenv" {
	t.Fatalf("expected base url from .env, got %q", cfg.BaseURL)
}
```

Add a second test that writes the same `.env` file, sets `MEMOS_URL` and `MEMOS_API_KEY` in the real environment, calls `LoadFromEnv()`, and asserts the environment values win.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/config -run TestLoadFromEnv`
Expected: FAIL because `LoadFromEnv()` currently only reads `os.Getenv()` and ignores `.env`.

**Step 3: Write minimal implementation**

Do not change tests in this step. Wait for the failure before touching production code.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/config -run TestLoadFromEnv`
Expected: PASS after Task 2 implementation.

**Step 5: Commit**

```bash
git add internal/config/config_test.go
git commit -m "test: cover dotenv config loading"
```

### Task 2: Implement dotenv-backed config loading

**Files:**
- Modify: `internal/config/config.go`
- Modify: `go.mod`
- Modify: `go.sum`
- Test: `internal/config/config_test.go`

**Step 1: Write the failing test**

Use the tests from Task 1 as the specification. No additional production behavior should be added before the failing tests are observed.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/config -run TestLoadFromEnv`
Expected: FAIL with empty config fields where `.env` values were expected.

**Step 3: Write minimal implementation**

Add `github.com/joho/godotenv` and update `LoadFromEnv()` to follow this exact pattern:

```go
func LoadFromEnv() Config {
	dotenvValues, err := godotenv.Read()
	if err != nil {
		dotenvValues = map[string]string{}
	}

	return Config{
		BaseURL:     loadValue("MEMOS_URL", dotenvValues),
		APIKey:      loadValue("MEMOS_API_KEY", dotenvValues),
		AdminAPIKey: loadValue("MEMOS_ADMIN_API_KEY", dotenvValues),
		DefaultTag:  loadValue("DEFAULT_TAG", dotenvValues),
	}
}

func loadValue(key string, dotenvValues map[string]string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return strings.TrimSpace(dotenvValues[key])
}
```

Keep behavior tolerant of missing `.env` by treating read errors as empty dotenv data.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/config -run TestLoadFromEnv`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/config/config.go internal/config/config_test.go go.mod go.sum
git commit -m "feat: load config from dotenv fallback"
```

### Task 3: Update docs and run full verification

**Files:**
- Modify: `README.md`

**Step 1: Write the failing test**

Documentation change only; no automated test needed.

**Step 2: Run test to verify it fails**

Not applicable for documentation.

**Step 3: Write minimal implementation**

Update the configuration section in `README.md` to explicitly state:
- The CLI automatically reads `.env` from the current working directory.
- System environment variables override values from `.env`.

Prefer changing only the existing configuration paragraph near the env example.

**Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS.

**Step 5: Commit**

```bash
git add README.md
git commit -m "docs: clarify dotenv precedence"
```
