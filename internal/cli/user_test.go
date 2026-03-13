package cli

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserListUsesAdminKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer admin-456" {
			t.Fatalf("expected admin bearer token, got %q", got)
		}
		_, _ = w.Write([]byte(`{"users":[{"name":"users/1","nickname":"alice"}]}`))
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")
	t.Setenv("MEMOS_ADMIN_API_KEY", "admin-456")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"user", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected user list to succeed, got %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("alice")) {
		t.Fatalf("expected user output to contain alice, got %s", buf.String())
	}
}

func TestUserListJSONOutputsUsersObject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"users":[{"name":"users/1","nickname":"alice"}]}`))
	}))
	defer server.Close()

	t.Setenv("MEMOS_URL", server.URL)
	t.Setenv("MEMOS_API_KEY", "token-123")
	t.Setenv("MEMOS_ADMIN_API_KEY", "admin-456")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--json", "user", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected user list json to succeed, got %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte(`"users"`)) {
		t.Fatalf("expected JSON users output, got %s", buf.String())
	}
}
