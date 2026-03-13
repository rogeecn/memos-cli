# Memo Detail Output Format Design

## Background

The CLI now renders memo list output as block-style text with ID and publish time. Single-memo detail output still uses the old two-line format from `WriteMemoDetail`, which makes `memo get` and other commands that return one memo look inconsistent with the updated list view.

## Goal

Unify single-memo detail output with the new list style by rendering the memo ID, local readable publish time, and content in the same order, without adding a separator line for single-item views.

## Chosen Approach

Update `output.WriteMemoDetail` to share the same display conventions as list output:

1. Show pure ID by trimming `memos/`.
2. Show formatted publish time on the next line when available.
3. Show trimmed content below it.
4. Omit the separator line for detail output.

This keeps detail output compact while visually aligned with the list renderer.

## Why This Approach

- Keeps formatting centralized in `internal/output`.
- Makes `memo get`, create/update responses, comment creation, and tag updates consistent.
- Avoids duplicating formatting logic per command.
- Preserves JSON output behavior entirely.

## Scope

- Modify `internal/output/memos.go`.
- Add or update CLI tests for `memo get` output.
- Leave list output unchanged.

## Testing Strategy

Use TDD:

1. Add a failing `memo get` text-output test with `createTime`.
2. Run the focused test and confirm it fails for the expected reason.
3. Implement the minimal change in `WriteMemoDetail`.
4. Re-run the focused test, then `go test ./...`.

## Success Criteria

- `memo get` outputs `> <id>` on the first line.
- Publish time appears on the second line in local readable format.
- Content appears below the timestamp.
- No separator line is added for detail output.
- Full Go test suite stays green.
