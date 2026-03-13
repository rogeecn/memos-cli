package cli

import (
	"bytes"
	"testing"
)

func TestConfigCheckShowsConfiguredValues(t *testing.T) {
	t.Setenv("MEMOS_URL", "https://memos.example.com")
	t.Setenv("MEMOS_API_KEY", "token-123")
	t.Setenv("MEMOS_ADMIN_API_KEY", "admin-456")
	t.Setenv("DEFAULT_TAG", "cli")

	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"config", "check"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected config check to succeed, got error: %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("MEMOS_URL")) {
		t.Fatalf("expected output to contain MEMOS_URL, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("configured")) {
		t.Fatalf("expected output to mention configured status, got: %s", output)
	}
}
