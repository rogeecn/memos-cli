package cli

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
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

	if !bytes.Contains(buf.Bytes(), []byte("memos/123")) {
		t.Fatalf("expected create output to contain memo name, got %s", buf.String())
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
