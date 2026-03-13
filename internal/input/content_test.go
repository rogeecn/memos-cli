package input

import (
	"strings"
	"testing"
)

func TestMergeDefaultTagAppendsHashTag(t *testing.T) {
	result := MergeTags("hello", []string{"cli"})
	if !strings.Contains(result, "#cli") {
		t.Fatalf("expected merged content to contain #cli, got %q", result)
	}
}

func TestRemoveTagDeletesHashTag(t *testing.T) {
	result := RemoveTag("hello #cli #work", "cli")
	if strings.Contains(result, "#cli") {
		t.Fatalf("expected #cli to be removed, got %q", result)
	}
}
