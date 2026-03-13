# Memo List Pagination Design

## Background

The Go CLI currently exposes `memo list`, but it only performs a single `ListMemos` request and prints the first response page. The API client in `internal/memos/client.go` supports only a `filter` query parameter, so the CLI cannot request subsequent pages or surface pagination metadata to users.

## Goal

Add explicit pagination support to `memo list` so users can request additional data with `--page-size` and `--page-token`, while preserving script-friendly JSON output and simple terminal output.

## Chosen Approach

Extend the Memos HTTP client request/response types to include pagination fields:

1. `ListMemosParams.PageSize`
2. `ListMemosParams.PageToken`
3. `ListMemosResponse.NextPageToken`

Then wire these fields into `memo list` flags. For text output, keep rendering only the current page of memos and append a short hint when `nextPageToken` is present. For `--json`, serialize the full response object so callers can consume both memo data and the next-page token.

## Why This Approach

- Keeps pagination logic centralized in the existing Memos client.
- Matches current CLI patterns: flags in Cobra, API plumbing in `internal/memos`, rendering in `internal/output` or command layer.
- Supports both interactive human usage and automation without introducing an interactive pager.
- Solves the immediate gap with minimal surface area and low regression risk.

## Alternatives Considered

### 1. `--all` auto-fetch every page

This is convenient for small datasets, but it can unexpectedly fetch large result sets and does not directly address the user request for loading more data page by page.

### 2. Interactive `--more` mode

This is closer to a GUI-style “load more” flow, but it complicates Cobra command behavior, test setup, and piping/JSON scenarios. It is better as a follow-up once the underlying pagination primitives exist.

## Scope

- Modify `internal/memos/client.go` to send pagination query parameters and decode `nextPageToken`.
- Modify `internal/cli/memo.go` to accept `--page-size` and `--page-token`.
- Add focused CLI tests for pagination query propagation and next-page handling.
- Update `README.md` Go CLI usage examples to mention pagination flags.

## Error Handling

- Invalid pagination values should rely on Cobra flag typing where possible.
- Empty `nextPageToken` means there is no additional page and should not print a hint.
- Existing API error handling remains unchanged.

## Testing Strategy

Use TDD for the CLI behavior:

1. Write a failing test proving `memo list --page-size --page-token` sends both query parameters.
2. Write a failing test proving text output includes a next-page hint when the API returns `nextPageToken`.
3. Write a failing test proving JSON output includes `nextPageToken`.
4. Run targeted tests to confirm they fail for the expected reason.
5. Implement the minimal production changes.
6. Re-run targeted tests, then run `go test ./...`.

## Success Criteria

- Users can request subsequent pages with `memo list --page-size N --page-token TOKEN`.
- Text output remains readable and tells users how to fetch the next page.
- JSON output includes pagination metadata for automation.
- Existing list/search/filter behavior continues to work.
