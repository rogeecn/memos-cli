# Memos Go CLI Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a Go-based CLI for Memos that mirrors the current MCP feature set and adds terminal-friendly workflows such as human-readable output, JSON output, stdin, and editor input.

**Architecture:** Keep the existing Python MCP implementation as-is for MCP use. Add a separate Go CLI inside this repository that talks directly to the Memos HTTP API, with clear package boundaries for config, API client, command wiring, output formatting, and content input helpers.

**Tech Stack:** Go, Cobra, standard library `net/http`, `encoding/json`, `testing`, `httptest`.

---

### Task 1: Create Go CLI skeleton
- Initialize `go.mod`
- Add `cmd/memos/main.go`
- Add `internal/cli/root.go`
- Add first smoke tests for `--help`

### Task 2: Add config loading
- Support env vars and global flags
- Add config validation helpers
- Add `config check` command with tests

### Task 3: Add Memos HTTP client
- Implement authenticated requests
- Add memo/user/comment API methods
- Add client tests with `httptest`

### Task 4: Add read commands
- Implement `memo list`, `memo get`, `search`, `filter`
- Add text and JSON output tests

### Task 5: Add write commands
- Implement `memo create`, `memo update`, `memo delete`
- Implement `comment create`, `tag remove`
- Add stdin and `--edit` content input helpers with tests

### Task 6: Update docs and verify
- Document Go CLI usage in `README.md`
- Run focused tests and `go test ./...`
- Summarize remaining gaps and next-step enhancements
