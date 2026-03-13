package cli

import (
	"bytes"
	"testing"
)

func TestRootCommandShowsHelp(t *testing.T) {
	cmd := NewRootCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected help command to succeed, got error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Fatal("expected help output, got empty string")
	}
	if !bytes.Contains([]byte(output), []byte("Memos")) {
		t.Fatalf("expected help output to mention Memos, got: %s", output)
	}
}
