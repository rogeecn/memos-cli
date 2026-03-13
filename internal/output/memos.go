package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/rogeecn/memos-cli/internal/memos"
)

func WriteMemoList(w io.Writer, items []memos.Memo) error {
	for _, item := range items {
		content := strings.TrimSpace(item.Content)
		if _, err := fmt.Fprintf(w, "%s\t%s\n", item.Name, content); err != nil {
			return err
		}
	}
	if len(items) == 0 {
		_, err := fmt.Fprintln(w, "No memos found.")
		return err
	}
	return nil
}

func WriteMemoDetail(w io.Writer, item memos.Memo) error {
	_, err := fmt.Fprintf(w, "%s\n%s\n", item.Name, strings.TrimSpace(item.Content))
	return err
}
