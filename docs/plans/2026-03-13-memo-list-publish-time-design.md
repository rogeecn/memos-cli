# Memo List Publish Time Design

## Background

The human-readable memo list output already renders each memo as a block with ID, content, and separator. The user now wants the publish time shown directly below the ID line. The current `memos.Memo` struct does not expose a timestamp field, so the output layer cannot render it yet.

## Goal

Show each memo's publish time on the line below the ID, formatted as local human-readable time, while preserving the existing block layout and leaving JSON output unchanged.

## Chosen Approach

Add a timestamp field to `memos.Memo` for the API's creation time, then format it in `output.WriteMemoList` using Go's local time conversion and a stable layout like `2006-01-02 15:04`. If parsing fails or the field is missing, print an empty timestamp line only if we explicitly choose to preserve layout; otherwise skip the time line. For this change, prefer rendering the parsed local time line when available and falling back to the raw value only if parsing is invalid.

## Why This Approach

- Keeps API data modeling in `internal/memos` and presentation in `internal/output`.
- Avoids adding command-specific formatting logic.
- Produces readable output for terminal users without changing structured JSON behavior.
- Maintains backward compatibility with existing API responses if a timestamp is absent.

## Alternatives Considered

### 1. Render raw API timestamp

Simpler, but less readable in terminal output.

### 2. Render relative time

Friendly, but less stable for tests and less precise for scripting or screenshots.

### 3. Render UTC formatted time

Stable, but less aligned with the user's explicit preference for local readable time.

## Scope

- Modify `internal/memos/client.go` to expose the memo creation/publish timestamp field.
- Modify `internal/output/memos.go` to render the timestamp below the ID.
- Update CLI tests in `internal/cli/memo_read_test.go` for the new block format.
- Keep JSON output untouched apart from including the mapped field if the API already returns it.

## Testing Strategy

Use TDD:

1. Add a failing list-output test with API `createTime` values.
2. Run the focused CLI test to verify the failure is due to missing timestamp rendering.
3. Implement the smallest client/output change.
4. Re-run the focused test, then `go test ./...`.

## Success Criteria

- List output prints ID on the first line.
- List output prints local readable publish time on the second line.
- Content stays below the timestamp.
- Block separator and blank line behavior remain unchanged.
- Full test suite remains green.
