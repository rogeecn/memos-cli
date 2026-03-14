# Memo Image Upload Design

## Background

The CLI currently creates memos by sending only `content` and `visibility` to `POST /api/v1/memos`. It supports Markdown text and tags, but it has no way to upload local image files or attach them to a memo. Official Memos API documentation now exposes attachment-based upload endpoints rather than the older resource model.

## Goal

Add image upload support to `memo create` so users can run `memos memo create "正文" --image a.png --image b.jpg` and get a memo created first, then have the provided local images uploaded and bound to that memo as attachments.

## Chosen Approach

Use a three-step flow based on the official attachment API:

1. Create the memo with the existing `POST /api/v1/memos` flow.
2. For each `--image` path, read the local file, infer its MIME type, base64-encode the bytes through Go's standard JSON handling for `[]byte`, and call `POST /api/v1/attachments`.
3. After all uploads succeed, call `PATCH /api/v1/memos/{memo}/attachments` once with the created attachments.

The CLI command surface remains backward-compatible. If `--image` is not provided, `memo create` behaves exactly as it does today.

## Why This Approach

- Matches the current official Memos API design using attachments.
- Keeps the memo creation path stable for text-only usage.
- Makes upload and binding failures visible at the correct step.
- Avoids depending on undocumented behavior from setting the `memo` field alone during upload.
- Lets the client model attachments explicitly for future commands if needed.

## Scope

- Modify `internal/cli/memo_write.go` to add repeatable `--image` flags and orchestrate create-upload-bind flow.
- Modify `internal/memos/client.go` to add attachment request and response models plus `CreateAttachment` and `SetMemoAttachments` methods.
- Add small file-handling helpers for reading image files and deriving MIME types.
- Add CLI and/or client tests covering the new HTTP interactions.
- Update README command examples to document `--image` usage.

## Non-Goals

- Adding image upload support to `memo update`.
- Adding generic non-image attachment workflows.
- Rewriting memo content to inject Markdown image links.
- Implementing rollback if memo creation succeeds but later upload/bind steps fail.

## Data Flow

1. Parse `memo create <content> --image <path>...`.
2. Build the existing memo payload and create the memo.
3. For each image path:
   - read the file from disk;
   - determine the filename from the path;
   - infer MIME type from file extension, falling back to `application/octet-stream`;
   - send an attachment payload containing `filename`, `type`, `content`, and `memo`.
4. Collect the returned attachments.
5. Send one set-attachments request for the created memo.
6. Print the created memo using the existing output path.

## Error Handling

- If memo creation fails, return the API error immediately.
- If an image path cannot be read, return an error that includes the path.
- If MIME detection fails, fall back to `application/octet-stream` rather than rejecting the file.
- If any attachment upload fails, stop immediately and return the error. The memo remains created.
- If attachment binding fails, return the error and keep the memo created state unchanged.

## Testing Strategy

Use TDD and verify the behavior end to end at the CLI level first:

1. Add a failing CLI test for `memo create --image` that asserts request order and payload shape.
2. Run the focused test and confirm it fails before implementation.
3. Implement minimal client and CLI changes.
4. Re-run the focused test until it passes.
5. Run `go test ./...` to verify the full suite.

## Success Criteria

- `memos memo create "text" --image a.png` creates the memo and uploads the image.
- Multiple `--image` flags upload multiple files and bind them in one final attachment request.
- Text-only `memo create` behavior remains unchanged.
- Tests verify attachment payloads and binding endpoint usage.
- README includes a discoverable image-upload example.
