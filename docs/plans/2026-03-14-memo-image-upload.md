# Memo Image Upload Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add `memo create --image` so the CLI can create a memo, upload local image files as attachments, and bind them to the new memo.

**Architecture:** Keep the existing memo-create command as the entrypoint and extend it with a repeatable `--image` flag. Model attachment API requests in `internal/memos/client.go`, then orchestrate the create-upload-bind sequence in `internal/cli/memo_write.go` using small file-loading helpers and focused `httptest` coverage.

**Tech Stack:** Go, Cobra, standard library `os`, `path/filepath`, `mime`, `net/http`, `encoding/json`, `testing`, `httptest`.

---

### Task 1: Add failing CLI upload flow test

**Files:**
- Modify: `internal/cli/memo_write_test.go`

**Step 1: Write the failing test**

Add `TestMemoCreateUploadsImagesAndBindsAttachments` using `httptest.NewServer` and a temporary PNG file. Make the server assert this request sequence:

1. `POST /api/v1/memos`
2. `POST /api/v1/attachments`
3. `PATCH /api/v1/memos/123/attachments`

Assert the attachment upload request body includes these JSON fields:

```json
{
  "filename": "pic.png",
  "type": "image/png",
  "content": "<base64 bytes>",
  "memo": "memos/123"
}
```

Assert the binding request body includes:

```json
{
  "name": "memos/123",
  "attachments": [
    {"name": "attachments/att-1", "filename": "pic.png", "type": "image/png"}
  ]
}
```

Run the command with arguments equivalent to:

```bash
memo create "hello" --image /tmp/pic.png
```

**Step 2: Run test to verify it fails**

Run: `TZ=UTC go test ./internal/cli -run '^TestMemoCreateUploadsImagesAndBindsAttachments$'`
Expected: FAIL because `memo create` does not yet define `--image` or call attachment endpoints.

**Step 3: Write minimal implementation**

Do not implement yet.

**Step 4: Run test to verify it passes**

Run after implementation: `TZ=UTC go test ./internal/cli -run '^TestMemoCreateUploadsImagesAndBindsAttachments$'`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 2: Add attachment client support

**Files:**
- Modify: `internal/memos/client.go`
- Modify: `internal/cli/memo_write_test.go`
- Test: `internal/memos/client_test.go`

**Step 1: Write the failing test**

If the CLI test from Task 1 does not already exercise the client payload sufficiently, add a client-level test that verifies:

- `CreateAttachment` sends `POST /api/v1/attachments`
- `SetMemoAttachments` sends `PATCH /api/v1/memos/123/attachments`
- JSON field names stay camelCase (`filename`, `type`, `content`, `memo`, `attachments`)

Use mock responses like:

```json
{"name":"attachments/att-1","filename":"pic.png","type":"image/png","size":"4","memo":"memos/123"}
```

and an empty `{}` body for set-attachments.

**Step 2: Run test to verify it fails**

Run: `TZ=UTC go test ./internal/memos -run 'Attachment|SetMemoAttachments'`
Expected: FAIL because the client types and methods do not exist yet.

**Step 3: Write minimal implementation**

Add attachment types and methods similar to:

```go
type Attachment struct {
    Name         string `json:"name,omitempty"`
    Filename     string `json:"filename,omitempty"`
    Content      []byte `json:"content,omitempty"`
    ExternalLink string `json:"externalLink,omitempty"`
    Type         string `json:"type,omitempty"`
    Size         string `json:"size,omitempty"`
    Memo         string `json:"memo,omitempty"`
}

type SetMemoAttachmentsPayload struct {
    Name        string       `json:"name"`
    Attachments []Attachment `json:"attachments"`
}
```

Add methods with these paths:

```go
func (c *Client) CreateAttachment(payload Attachment) (Attachment, error)
func (c *Client) SetMemoAttachments(memoID string, payload SetMemoAttachmentsPayload) error
```

Use the existing `doJSON` helper for both.

**Step 4: Run test to verify it passes**

Run: `TZ=UTC go test ./internal/memos -run 'Attachment|SetMemoAttachments'`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 3: Implement CLI image upload orchestration

**Files:**
- Modify: `internal/cli/memo_write.go`
- Modify: `internal/cli/memo_write_test.go`

**Step 1: Write the failing test**

Reuse `TestMemoCreateUploadsImagesAndBindsAttachments` from Task 1.

**Step 2: Run test to verify it fails**

Run: `TZ=UTC go test ./internal/cli -run '^TestMemoCreateUploadsImagesAndBindsAttachments$'`
Expected: FAIL because the command does not yet parse image flags or upload files.

**Step 3: Write minimal implementation**

Update `newMemoCreateCommand` to:

- declare `var images []string`
- register `cmd.Flags().StringSliceVar(&images, "image", nil, "Local image paths to upload")`
- create the memo first using existing logic
- if `len(images) == 0`, return the current output immediately
- otherwise, for each image path:
  - `os.ReadFile(path)`
  - `filepath.Base(path)` for `filename`
  - `mime.TypeByExtension(filepath.Ext(path))`, fallback to `application/octet-stream`
  - call `client.CreateAttachment(...)` with `Memo: memo.Name`
- call `client.SetMemoAttachments(...)` once with all returned attachments
- preserve existing text and JSON output behavior

Return descriptive errors that include the path for local file read failures.

**Step 4: Run test to verify it passes**

Run: `TZ=UTC go test ./internal/cli -run '^TestMemoCreateUploadsImagesAndBindsAttachments$'`
Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 4: Document image upload usage

**Files:**
- Modify: `README.md`

**Step 1: Write the failing test**

No automated test required.

**Step 2: Run test to verify it fails**

Skip.

**Step 3: Write minimal implementation**

Add a new example under memo write operations:

```bash
memos memo create "旅行记录" --image cover.png
memos memo create "相册更新" --image a.jpg --image b.jpg
```

Keep the docs scoped to `memo create` only.

**Step 4: Run test to verify it passes**

No automated test required.

**Step 5: Commit**

Do not commit unless explicitly requested.

### Task 5: Verify focused and full tests

**Files:**
- No further file changes required unless verification reveals regressions.

**Step 1: Write the failing test**

No additional test required.

**Step 2: Run test to verify it fails**

Skip.

**Step 3: Write minimal implementation**

None.

**Step 4: Run test to verify it passes**

Run these commands:

```bash
TZ=UTC go test ./internal/cli -run '^TestMemo(CreateUploadsImagesAndBindsAttachments|CreateAppendsDefaultTag)$'
TZ=UTC go test ./internal/memos -run 'Attachment|SetMemoAttachments'
TZ=UTC go test ./...
```

Expected: PASS.

**Step 5: Commit**

Do not commit unless explicitly requested.
