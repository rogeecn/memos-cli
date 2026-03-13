package memos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListMemosSendsBearerAndFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memos" {
			t.Fatalf("expected path /api/v1/memos, got %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token-123" {
			t.Fatalf("expected bearer token, got %q", got)
		}
		if got := r.URL.Query().Get("filter"); got != "content.contains('cli')" {
			t.Fatalf("expected filter query, got %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"memos": []map[string]any{{"name": "memos/123", "content": "hello"}}})
	}))
	defer server.Close()

	client := NewClient(server.URL, "token-123", "")
	resp, err := client.ListMemos(ListMemosParams{Filter: "content.contains('cli')"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Memos) != 1 {
		t.Fatalf("expected one memo, got %d", len(resp.Memos))
	}
}

func TestListUsersRequiresAdminKey(t *testing.T) {
	client := NewClient("https://memos.example.com", "token-123", "")

	_, err := client.ListUsers()
	if err == nil {
		t.Fatal("expected admin key error, got nil")
	}
}

func TestGetMemoUsesMemoNamePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memos/123" {
			t.Fatalf("expected memo path, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"name": "memos/123", "content": "memo detail"})
	}))
	defer server.Close()

	client := NewClient(server.URL, "token-123", "")
	memo, err := client.GetMemo("123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if memo.Content != "memo detail" {
		t.Fatalf("expected memo detail content, got %q", memo.Content)
	}
}

func TestCreateCommentUsesMemoNamePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memos/123/comments" {
			t.Fatalf("expected memo comment path, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"name": "memos/comments/1"})
	}))
	defer server.Close()

	client := NewClient(server.URL, "token-123", "")
	_, err := client.CreateComment("123", MemoPayload{Content: "note", Visibility: "PRIVATE"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
