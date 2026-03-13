# Memo Detail Output Format Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make single-memo detail output match the new block-style memo formatting by showing pure ID, publish time, and content in that order.

**Architecture:** Reuse the existing `Memo` data shape and time-formatting helper in `internal/output/memos.go`, so list and detail views stay visually consistent without moving formatting logic into CLI commands. Verify the behavior through a CLI test that exercises `memo get` end to end.

**Tech Stack:** Go, standard library `fmt`, `strings`, `time`, `testing`, `httptest`.

---

### Task 1: Add failing memo-detail format test

**Files:**
- Modify: `internal/cli/memo_read_test.go`

**Step 1: Write the failing test**

Add a focused `memo get` test whose mock API response includes:

```json
{"name":"memos/123","content":"ship cli","createTime":"2026-03-13T09:30:00Z"}
```

With `TZ=UTC`, assert the exact output is:

```text
> 123
2026-03-13 09:30
ship cli
```

**Step 2: Run test to verify it fails**

Run: `TZ=UTC go test ./internal/cli -run '^TestMemoGetPrintsBlockDetailFormat$'`
Expected: FAIL because detail output still uses the old `name + content` format.

**Step 3: Write minimal implementation**

Do not implement yet.

**Step 4: Run test to verify it passes**

Run after implementation: `TZ=UTC go test ./internal/cli -run '^TestMemoGetPrintsBlockDetailFormat$'`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 2: Implement detail formatter update

**Files:**
- Modify: `internal/output/memos.go`
- Modify: `internal/cli/memo_read_test.go`

**Step 1: Write the failing test**

Reuse the failing test from Task 1.

**Step 2: Run test to verify it fails**

Run: `TZ=UTC go test ./internal/cli -run '^TestMemoGetPrintsBlockDetailFormat$'`
Expected: FAIL because detail formatting has not been updated yet.

**Step 3: Write minimal implementation**

Update `WriteMemoDetail` to:
- trim `memos/` prefix for display ID
- render local readable time using the existing create-time formatter when present
- print content after the timestamp
- omit any separator line

**Step 4: Run test to verify it passes**

Run: `TZ=UTC go test ./internal/cli -run '^TestMemoGetPrintsBlockDetailFormat$'`
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
