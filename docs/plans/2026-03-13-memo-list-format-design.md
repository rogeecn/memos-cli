# Memo List Output Format Design

## Background

The Go CLI currently renders memo lists in a single-line tab-separated format from `internal/output/memos.go`. The user wants list output to render each memo as a small block with a quoted ID line, content line, and separator line, with a blank line between adjacent memo blocks.

## Goal

Change human-readable memo list output to this structure:

```text
> 123
memo content
-----------------
```

with one blank line between blocks, while preserving JSON output behavior.

## Chosen Approach

Update `output.WriteMemoList` so all commands that rely on human-readable memo list rendering (`memo list`, `search`, `filter`) automatically inherit the new format. Extract the display ID from `Memo.Name` by trimming the `memos/` prefix when present.

## Why This Approach

- Keeps formatting logic centralized in the output layer.
- Avoids duplicating rendering logic across CLI commands.
- Preserves existing command behavior apart from the requested display format.
- Keeps JSON output untouched.

## Scope

- Modify `internal/output/memos.go` to render block output.
- Add or update tests in `internal/cli/memo_read_test.go` to cover the new text format.
- Leave JSON output and detail output unchanged.

## Testing Strategy

Use TDD:

1. Add a failing test for the new block format in list output.
2. Run the focused CLI test to verify it fails for the expected reason.
3. Implement the smallest change in `internal/output/memos.go`.
4. Re-run the focused test, then `go test ./...`.

## Success Criteria

- Each listed memo renders as `> <id>` then content then `-----------------`.
- Adjacent memo blocks are separated by one blank line.
- IDs are shown without the `memos/` prefix.
- JSON output remains unchanged.
