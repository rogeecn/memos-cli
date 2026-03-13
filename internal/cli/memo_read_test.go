package cli

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMemoCommandHelpExists(t *testing.T) {
	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"memo", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected memo help to succeed, got %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("list")) {
		t.Fatalf("expected memo help to mention list subcommand, got %s", buf.String())
	}
}

func TestMemoListPrintsMemoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memos" {
			t.Fatalf("expected memos path, got %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"memos":[{"name":"memos/123","content":"ship cli","visibility":"PRIVATE"}]}`))
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"memo", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected memo list to succeed, got %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("ship cli")) {
		t.Fatalf("expected memo content in output, got %s", buf.String())
	}
}

func TestSearchUsesFilterExpression(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("filter"); got != "content.contains('roadmap')" {
			t.Fatalf("expected search filter, got %q", got)
		}
		_, _ = w.Write([]byte(`{"memos":[{"name":"memos/456","content":"roadmap notes"}]}`))
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"search", "roadmap"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected search command to succeed, got %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("roadmap notes")) {
		t.Fatalf("expected search output to contain memo content, got %s", buf.String())
	}
}

func TestFilterUsesExplicitExpression(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("filter"); got != "visibility == 'PRIVATE'" {
			t.Fatalf("expected explicit filter, got %q", got)
		}
		_, _ = w.Write([]byte(`{"memos":[{"name":"memos/789","content":"private memo"}]}`))
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"filter", "--expr", "visibility == 'PRIVATE'"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected filter command to succeed, got %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("private memo")) {
		t.Fatalf("expected filter output to contain memo content, got %s", buf.String())
	}
}
