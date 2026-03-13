# README Refresh Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Rewrite `README.md` so it accurately documents this repository as a Go CLI for common Memos operations.

**Architecture:** Replace the obsolete Python/MCP-oriented README with a CLI-first document organized around real command groups, current configuration loading behavior, and verified examples. Validate content against the Cobra command tree and current config code instead of relying on stale prose.

**Tech Stack:** Markdown, Go, Cobra CLI.

---

### Task 1: Verify the live command surface

**Files:**
- Read: `main.go`
- Read: `internal/cli/root.go`
- Read: `internal/cli/config.go`
- Read: `internal/cli/memo.go`
- Read: `internal/cli/memo_write.go`
- Read: `internal/cli/comment.go`
- Read: `internal/cli/tag.go`
- Read: `internal/cli/user.go`

**Step 1: Write the failing test**

No code test required. This is a documentation verification task.

**Step 2: Run test to verify it fails**

Skip.

**Step 3: Write minimal implementation**

Extract the exact command names, subcommands, required flags, and notable behaviors that must appear in the README:
- binary name `memos`
- global `--json`
- `config check`
- `memo list|get|create|update|delete`
- `search`
- `filter --expr`
- `comment create`
- `tag remove`
- `user list`
- `memo delete --yes`
- `user list` requiring admin credentials

**Step 4: Run test to verify it passes**

Run: `go run . --help`
Expected: top-level command help prints and matches the documented binary/command tree.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 2: Rewrite README content

**Files:**
- Modify: `README.md`

**Step 1: Write the failing test**

Identify the current mismatches in the existing README:
- wrong project identity (Python MCP server)
- wrong runtime assumptions (`uv`, Python)
- stale feature lists and examples
- stale source layout references

**Step 2: Run test to verify it fails**

Manual verification: read `README.md` and confirm it does not match the current Go CLI. Expected: FAIL by inspection.

**Step 3: Write minimal implementation**

Replace the README with a complete CLI-focused structure that includes:
- concise project introduction
- current feature overview
- installation via `go install`
- configuration variables and `.env` precedence
- quick start section
- command reference with real syntax
- JSON output and pagination notes
- usage examples
- development and test notes
- license

**Step 4: Run test to verify it passes**

Read the rewritten `README.md` and confirm every documented command and flag maps to the current code.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 3: Validate command examples

**Files:**
- Modify: `README.md` if validation reveals mismatches

**Step 1: Write the failing test**

No new tests required.

**Step 2: Run test to verify it fails**

Run these commands:
- `go run . --help`
- `go run . memo --help`
- `go run . comment --help`
- `go run . tag --help`
- `go run . user --help`
Expected: help output confirms the documented command groups and names.

**Step 3: Write minimal implementation**

If help output reveals any wording or syntax mismatch in the README, correct only those lines.

**Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS, demonstrating the repo still builds and tests cleanly after the documentation update.

**Step 5: Commit**

Do not commit unless explicitly requested.
