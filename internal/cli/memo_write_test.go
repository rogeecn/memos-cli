package cli

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestMemoCreateAppendsDefaultTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		if !bytes.Contains(body, []byte("#cli")) {
			t.Fatalf("expected request body to include default tag, got %s", string(body))
		}
		_, _ = w.Write([]byte(`{"name":"memos/123","content":"hello\n #cli"}`))
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")
	t.Setenv("DEFAULT_TAG", "cli")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"memo", "create", "hello"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected memo create to succeed, got %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("> 123")) {
		t.Fatalf("expected create output to contain formatted memo id, got %s", buf.String())
	}
}

func TestMemoCreateUploadsImagesAndBindsAttachments(t *testing.T) {
	imagePath := filepath.Join(t.TempDir(), "pic.png")
	imageBytes := []byte{0x89, 0x50, 0x4e, 0x47}
	if err := os.WriteFile(imagePath, imageBytes, 0o600); err != nil {
		t.Fatalf("write temp image: %v", err)
	}

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch callCount {
		case 1:
			if r.Method != http.MethodPost || r.URL.Path != "/api/v1/memos" {
				t.Fatalf("expected first call POST /api/v1/memos, got %s %s", r.Method, r.URL.Path)
			}
			_, _ = w.Write([]byte(`{"name":"memos/123","content":"hello"}`))
		case 2:
			if r.Method != http.MethodPost || r.URL.Path != "/api/v1/attachments" {
				t.Fatalf("expected second call POST /api/v1/attachments, got %s %s", r.Method, r.URL.Path)
			}
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode attachment payload: %v", err)
			}
			if got := payload["filename"]; got != "pic.png" {
				t.Fatalf("expected filename pic.png, got %#v", got)
			}
			if got := payload["type"]; got != "image/png" {
				t.Fatalf("expected type image/png, got %#v", got)
			}
			if got := payload["memo"]; got != "memos/123" {
				t.Fatalf("expected memo memos/123, got %#v", got)
			}
			expectedContent := base64.StdEncoding.EncodeToString(imageBytes)
			if got := payload["content"]; got != expectedContent {
				t.Fatalf("expected base64 content %q, got %#v", expectedContent, got)
			}
			_, _ = w.Write([]byte(`{"name":"attachments/att-1","filename":"pic.png","type":"image/png","memo":"memos/123"}`))
		case 3:
			if r.Method != http.MethodPatch || r.URL.Path != "/api/v1/memos/123/attachments" {
				t.Fatalf("expected third call PATCH /api/v1/memos/123/attachments, got %s %s", r.Method, r.URL.Path)
			}
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode set attachments payload: %v", err)
			}
			if got := payload["name"]; got != "memos/123" {
				t.Fatalf("expected memo name memos/123, got %#v", got)
			}
			attachments, ok := payload["attachments"].([]any)
			if !ok || len(attachments) != 1 {
				t.Fatalf("expected one attachment, got %#v", payload["attachments"])
			}
			attachment, ok := attachments[0].(map[string]any)
			if !ok {
				t.Fatalf("expected attachment object, got %#v", attachments[0])
			}
			if got := attachment["name"]; got != "attachments/att-1" {
				t.Fatalf("expected attachment name attachments/att-1, got %#v", got)
			}
			if got := attachment["filename"]; got != "pic.png" {
				t.Fatalf("expected attachment filename pic.png, got %#v", got)
			}
			if got := attachment["type"]; got != "image/png" {
				t.Fatalf("expected attachment type image/png, got %#v", got)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		default:
			t.Fatalf("unexpected extra request %d", callCount)
		}
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"memo", "create", "hello", "--image", imagePath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected memo create with image to succeed, got %v", err)
	}

	if callCount != 3 {
		t.Fatalf("expected 3 requests, got %d", callCount)
	}
}

func TestMemoDeleteRequiresYesFlag(t *testing.T) {
	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"memo", "delete", "123"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected delete without --yes to fail")
	}
}

func TestCommentCreateUsesMemoPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memos/123/comments" {
			t.Fatalf("expected comment path, got %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"name":"memos/comments/1","content":"reply"}`))
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"comment", "create", "123", "reply"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected comment create to succeed, got %v", err)
	}
}

func TestTagRemoveUpdatesMemoContent(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch callCount {
		case 1:
			if r.Method != http.MethodGet {
				t.Fatalf("expected first call GET, got %s", r.Method)
			}
			_, _ = w.Write([]byte(`{"name":"memos/123","content":"hello #cli #work"}`))
		case 2:
			if r.Method != http.MethodPatch {
				t.Fatalf("expected second call PATCH, got %s", r.Method)
			}
			body, _ := io.ReadAll(r.Body)
			if bytes.Contains(body, []byte("#cli")) {
				t.Fatalf("expected removed tag to be absent, got %s", string(body))
			}
			_, _ = w.Write([]byte(`{"name":"memos/123","content":"hello #work"}`))
		default:
			t.Fatalf("unexpected extra request %d", callCount)
		}
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"tag", "remove", "123", "cli"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected tag remove to succeed, got %v", err)
	}
}

func TestMemoUpdateSendsPatchedContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("expected PATCH, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		if !bytes.Contains(body, []byte("updated text")) {
			t.Fatalf("expected request body to include updated text, got %s", string(body))
		}
		_, _ = w.Write([]byte(`{"name":"memos/123","content":"updated text"}`))
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"memo", "update", "123", "--content", "updated text"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected memo update to succeed, got %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("updated text")) {
		t.Fatalf("expected update output to contain updated text, got %s", buf.String())
	}
}

func TestMemoListJSONOutputsRawJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"memos":[{"name":"memos/123","content":"json memo"}]}`))
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--json", "memo", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected memo list --json to succeed, got %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte(`"memos"`)) {
		t.Fatalf("expected JSON output, got %s", buf.String())
	}
}
