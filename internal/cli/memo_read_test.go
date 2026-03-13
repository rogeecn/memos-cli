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

func TestMemoListPrintsBlockFormat(t *testing.T) {
	t.Setenv("TZ", "UTC")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memos" {
			t.Fatalf("expected memos path, got %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"memos":[{"name":"memos/123","content":"first memo","visibility":"PRIVATE","createTime":"2026-03-13T09:30:00Z"},{"name":"memos/456","content":"second memo","visibility":"PRIVATE","createTime":"2026-03-13T10:45:00Z"}]}`))
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

	expected := "> 123\n2026-03-13 09:30\nfirst memo\n-----------------\n\n> 456\n2026-03-13 10:45\nsecond memo\n-----------------\n"
	if buf.String() != expected {
		t.Fatalf("expected block format output %q, got %q", expected, buf.String())
	}
}

func TestMemoGetPrintsBlockDetailFormat(t *testing.T) {
	t.Setenv("TZ", "UTC")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memos/123" {
			t.Fatalf("expected memo get path, got %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"name":"memos/123","content":"ship cli","createTime":"2026-03-13T09:30:00Z"}`))
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"memo", "get", "123"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected memo get to succeed, got %v", err)
	}

	expected := "> 123\n2026-03-13 09:30\nship cli\n"
	if buf.String() != expected {
		t.Fatalf("expected block detail output %q, got %q", expected, buf.String())
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

func TestMemoListPaginationSendsPageParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("pageSize"); got != "20" {
			t.Fatalf("expected pageSize query param, got %q", got)
		}
		if got := r.URL.Query().Get("pageToken"); got != "cursor-2" {
			t.Fatalf("expected pageToken query param, got %q", got)
		}
		_, _ = w.Write([]byte(`{"memos":[]}`))
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"memo", "list", "--page-size", "20", "--page-token", "cursor-2"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected memo list pagination to succeed, got %v", err)
	}
}

func TestMemoListTextOutputIncludesNextPageTokenHint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"memos":[{"name":"memos/123","content":"ship cli","visibility":"PRIVATE"}],"nextPageToken":"cursor-2"}`))
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

	if !bytes.Contains(buf.Bytes(), []byte("Next page token: cursor-2")) {
		t.Fatalf("expected next page token hint in output, got %s", buf.String())
	}
}

func TestMemoListJSONOutputIncludesNextPageToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"memos":[{"name":"memos/123","content":"ship cli"}],"nextPageToken":"cursor-2"}`))
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

	if !bytes.Contains(buf.Bytes(), []byte(`"nextPageToken": "cursor-2"`)) {
		t.Fatalf("expected JSON output to include nextPageToken, got %s", buf.String())
	}
}
