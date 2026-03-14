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

func TestCreateAttachmentUsesAttachmentPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/attachments" {
			t.Fatalf("expected attachment path, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
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
		_ = json.NewEncoder(w).Encode(map[string]any{"name": "attachments/att-1", "filename": "pic.png", "type": "image/png", "memo": "memos/123"})
	}))
	defer server.Close()

	client := NewClient(server.URL, "token-123", "")
	attachment, err := client.CreateAttachment(Attachment{Filename: "pic.png", Content: []byte("png"), Type: "image/png", Memo: "memos/123"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if attachment.Name != "attachments/att-1" {
		t.Fatalf("expected attachment name, got %q", attachment.Name)
	}
}

func TestSetMemoAttachmentsUsesMemoAttachmentPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memos/123/attachments" {
			t.Fatalf("expected memo attachments path, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPatch {
			t.Fatalf("expected PATCH, got %s", r.Method)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["name"]; got != "memos/123" {
			t.Fatalf("expected name memos/123, got %#v", got)
		}
		attachments, ok := payload["attachments"].([]any)
		if !ok || len(attachments) != 1 {
			t.Fatalf("expected one attachment, got %#v", payload["attachments"])
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "token-123", "")
	err := client.SetMemoAttachments("123", SetMemoAttachmentsPayload{
		Name: "memos/123",
		Attachments: []Attachment{{
			Name:     "attachments/att-1",
			Filename: "pic.png",
			Type:     "image/png",
		}},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
