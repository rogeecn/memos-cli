# README Refresh Design

## Background

The current `README.md` still describes an older Python-based MCP server. The repository now primarily contains a Go CLI named `memos`, with Cobra commands under `internal/cli` and a Go entrypoint in `main.go`. This mismatch makes the project landing page misleading for both users and contributors.

## Goal

Replace the outdated README with a complete, Go-CLI-first document that accurately describes what this repository is today: a command-line client for common Memos API operations.

## Chosen Approach

Rewrite `README.md` from top to bottom as a user-facing CLI document. Keep the document focused on practical usage: installation, configuration, command overview, examples, JSON output, and development workflow. Do not preserve the old Python/MCP narrative except implicitly by removing it.

## Why This Approach

- Fixes the root problem instead of patching scattered sections in an obsolete document.
- Gives new users a trustworthy landing page aligned with the actual executable and command surface.
- Keeps scope tight by updating only `README.md`, per user direction.
- Avoids introducing migration-history noise that does not help current users.

## Scope

- Rewrite `README.md` to describe the Go CLI.
- Remove outdated references to Python, MCP server mode, `uv`, `streamable-http`, MCP resources/tools/prompts, and stale file paths.
- Document the currently implemented commands and configuration variables.
- Keep changes limited to `README.md`; do not update `.env.example` or other docs in this task.

## Proposed README Structure

1. Project overview
2. Feature overview
3. Installation
4. Configuration and `.env` loading behavior
5. Quick start
6. Command reference by area (`config`, `memo`, `search`, `filter`, `comment`, `tag`, `user`)
7. JSON output and pagination notes
8. Usage examples
9. Development notes and test command
10. License

## Content Constraints

- Use the real binary/command name exposed by Cobra: `memos`.
- Keep examples consistent with actual implemented flags and subcommands.
- Mention that `user list` requires `MEMOS_ADMIN_API_KEY`.
- Mention that `memo delete` requires `--yes`.
- Mention that environment variables override `.env` values when both exist.
- Avoid claiming features that are not implemented, such as editor mode or stdin-based memo creation.

## Testing Strategy

Because this is a documentation-only change, validate by:

1. Reading the command definitions under `internal/cli` to ensure every documented command exists.
2. Running targeted CLI help commands to confirm top-level and nested command names.
3. Reviewing the rewritten README for consistency with current code.

## Success Criteria

- `README.md` no longer describes the repository as a Python MCP server.
- Installation, configuration, and command examples reflect the Go CLI accurately.
- A new user can understand how to configure and use `memos` without consulting source code first.
- The document stays within the requested scope of updating only the README.
