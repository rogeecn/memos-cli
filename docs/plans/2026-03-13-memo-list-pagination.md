# Memo List Pagination Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add page-by-page pagination support to `memo list` so the CLI can fetch additional memo pages with explicit flags and expose `nextPageToken` in both human and JSON output.

**Architecture:** Extend the existing Memos API client request and response types with pagination fields, then pass them through the `memo list` Cobra command. Keep text rendering page-scoped and append a next-step hint only when another page exists. Preserve JSON mode by serializing the full API response object.

**Tech Stack:** Go, Cobra, standard library `net/url`, `encoding/json`, `testing`, `httptest`.

---

### Task 1: Add failing pagination CLI tests

**Files:**
- Modify: `internal/cli/memo_read_test.go`
- Test: `internal/cli/memo_read_test.go`

**Step 1: Write the failing test**

Add a test that runs:

```go
cmd.SetArgs([]string{"memo", "list", "--page-size", "20", "--page-token", "cursor-2"})
```

and asserts the handler receives:

```go
r.URL.Query().Get("pageSize") == "20"
r.URL.Query().Get("pageToken") == "cursor-2"
```

Add a second test whose API response contains:

```json
{"memos":[{"name":"memos/123","content":"ship cli"}],"nextPageToken":"cursor-2"}
```

and asserts text output contains `cursor-2`.

Add a third test that runs:

```go
cmd.SetArgs([]string{"--json", "memo", "list"})
```

against the same response and asserts output contains `"nextPageToken": "cursor-2"`.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run 'TestMemoList'`
Expected: FAIL because pagination flags or output are not implemented yet.

**Step 3: Write minimal implementation**

Do not implement yet; only proceed after observing the failure.

**Step 4: Run test to verify it passes**

Run after implementation: `go test ./internal/cli -run 'TestMemoList'`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 2: Add pagination fields to the Memos client

**Files:**
- Modify: `internal/memos/client.go`
- Test: `internal/cli/memo_read_test.go`

**Step 1: Write the failing test**

Reuse the CLI request-assertion test from Task 1 as the failing signal for missing query propagation.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run TestMemoListPassesPaginationQuery`
Expected: FAIL because `pageSize` and `pageToken` are absent from the request.

**Step 3: Write minimal implementation**

Extend the request/response types:

```go
type ListMemosParams struct {
    Filter    string
    PageSize  int
    PageToken string
}

type ListMemosResponse struct {
    Memos         []Memo `json:"memos"`
    NextPageToken string `json:"nextPageToken,omitempty"`
}
```

Encode query params only when values are present:

```go
if params.PageSize > 0 {
    query.Set("pageSize", strconv.Itoa(params.PageSize))
}
if strings.TrimSpace(params.PageToken) != "" {
    query.Set("pageToken", params.PageToken)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli -run TestMemoListPassesPaginationQuery`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 3: Wire `memo list` flags and next-page hint

**Files:**
- Modify: `internal/cli/memo.go`
- Modify: `internal/cli/memo_read_test.go`

**Step 1: Write the failing test**

Reuse the text-output and JSON-output tests from Task 1.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run 'TestMemoList(PrintsNextPageHint|JSONIncludesNextPageToken)'`
Expected: FAIL because command flags and hint logic do not exist yet.

**Step 3: Write minimal implementation**

Add flags to `memo list`:

```go
var pageSize int
var pageToken string
cmd.Flags().IntVar(&pageSize, "page-size", 0, "Number of memos to fetch")
cmd.Flags().StringVar(&pageToken, "page-token", "", "Token for the next memo page")
```

Pass them into `ListMemos`:

```go
response, err := client.ListMemos(memos.ListMemosParams{
    PageSize:  pageSize,
    PageToken: pageToken,
})
```

For text output, print memos first, then append a short hint if `response.NextPageToken != ""`.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli -run 'TestMemoList'`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 4: Update docs and verify the full suite

**Files:**
- Modify: `README.md`

**Step 1: Write the failing test**

No automated README test required.

**Step 2: Run test to verify it fails**

Skip; documentation-only change.

**Step 3: Write minimal implementation**

Update the Go CLI section to mention:

```bash
go run . memo list --page-size 20 --page-token <token>
```

and explain that `--json` includes `nextPageToken`.

**Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.
