# Memo List Publish Time Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a local readable publish-time line under each memo ID in human-readable list output.

**Architecture:** Extend the memo API model with the creation timestamp, then keep all formatting in `internal/output/memos.go`. Tests should exercise the CLI surface with `httptest` responses that include `createTime`, so we verify the full request/decode/render path.

**Tech Stack:** Go, standard library `time`, `fmt`, `strings`, `testing`, `httptest`.

---

### Task 1: Add failing publish-time format test

**Files:**
- Modify: `internal/cli/memo_read_test.go`

**Step 1: Write the failing test**

Add a focused test that returns two memos with `createTime` values like `2026-03-13T09:30:00Z` and asserts text output becomes:

```text
> 123
2026-03-13 09:30
first memo
-----------------

> 456
2026-03-13 10:45
second memo
-----------------
```

Use UTC in the test process so local-readable output is deterministic.

**Step 2: Run test to verify it fails**

Run: `TZ=UTC go test ./internal/cli -run TestMemoListPrintsBlockFormat`
Expected: FAIL because the timestamp line is not rendered yet.

**Step 3: Write minimal implementation**

Do not implement yet.

**Step 4: Run test to verify it passes**

Run after implementation: `TZ=UTC go test ./internal/cli -run TestMemoListPrintsBlockFormat`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 2: Add memo timestamp field and render it

**Files:**
- Modify: `internal/memos/client.go`
- Modify: `internal/output/memos.go`
- Modify: `internal/cli/memo_read_test.go`

**Step 1: Write the failing test**

Reuse the failing CLI test from Task 1.

**Step 2: Run test to verify it fails**

Run: `TZ=UTC go test ./internal/cli -run TestMemoListPrintsBlockFormat`
Expected: FAIL because the timestamp field is not decoded/rendered.

**Step 3: Write minimal implementation**

- Add `CreateTime string \`json:"createTime,omitempty"\`` to `memos.Memo`.
- In `WriteMemoList`, parse `CreateTime` with RFC3339 when present.
- Convert to local time and render it with layout `2006-01-02 15:04` on the line below the ID.
- If parsing fails, fall back to the raw `CreateTime` string so output still includes the publish time.

**Step 4: Run test to verify it passes**

Run: `TZ=UTC go test ./internal/cli -run TestMemoListPrintsBlockFormat`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 3: Verify full suite

**Files:**
- No further file changes required unless verification reveals regressions.

**Step 1: Write the failing test**

No additional test required.

**Step 2: Run test to verify it fails**

Skip.

**Step 3: Write minimal implementation**

None.

**Step 4: Run test to verify it passes**

Run: `TZ=UTC go test ./...`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.
