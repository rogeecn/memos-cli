# Memo List Output Format Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Update human-readable memo list output to render each memo as a block with a quoted ID, content, separator line, and a blank line between blocks.

**Architecture:** Keep formatting centralized in `internal/output/memos.go`, where all list-style memo commands already converge. Validate the behavior through CLI tests that exercise the existing command surface rather than testing formatting in isolation only.

**Tech Stack:** Go, standard library `fmt`, `strings`, `testing`, `httptest`.

---

### Task 1: Add failing output-format test

**Files:**
- Modify: `internal/cli/memo_read_test.go`

**Step 1: Write the failing test**

Add a test that returns two memos from the mock server and asserts the text output contains:

```text
> 123
first memo
-----------------

> 456
second memo
-----------------
```

This proves the output uses pure IDs, block separators, and a blank line between blocks.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run TestMemoListPrintsBlockFormat`
Expected: FAIL because output is still tab-separated.

**Step 3: Write minimal implementation**

Do not implement yet.

**Step 4: Run test to verify it passes**

Run after implementation: `go test ./internal/cli -run TestMemoListPrintsBlockFormat`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 2: Implement block rendering

**Files:**
- Modify: `internal/output/memos.go`
- Modify: `internal/cli/memo_read_test.go`

**Step 1: Write the failing test**

Reuse the failing CLI test from Task 1.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run TestMemoListPrintsBlockFormat`
Expected: FAIL because rendering has not changed yet.

**Step 3: Write minimal implementation**

Update `WriteMemoList` to:
- derive the display ID by trimming `memos/` from `item.Name`
- print `> <id>`
- print trimmed content on the next line
- print `-----------------`
- print one blank line between items, but not an extra leading block separator before the first item

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli -run TestMemoListPrintsBlockFormat`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 3: Verify full suite

**Files:**
- No further file changes required unless failures reveal regressions.

**Step 1: Write the failing test**

No additional test authoring required.

**Step 2: Run test to verify it fails**

Skip.

**Step 3: Write minimal implementation**

None.

**Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.
