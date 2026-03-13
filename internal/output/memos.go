package output

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rogeecn/memos-cli/internal/memos"
)

func WriteMemoList(w io.Writer, items []memos.Memo) error {
	if len(items) == 0 {
		_, err := fmt.Fprintln(w, "No memos found.")
		return err
	}

	for index, item := range items {
		content := strings.TrimSpace(item.Content)
		id := strings.TrimPrefix(item.Name, "memos/")
		timeLine := formatMemoCreateTime(item.CreateTime)
		if timeLine != "" {
			if _, err := fmt.Fprintf(w, "> %s\n%s\n%s\n-----------------\n", id, timeLine, content); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(w, "> %s\n%s\n-----------------\n", id, content); err != nil {
				return err
			}
		}
		if index < len(items)-1 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
	}

	return nil
}

func WriteMemoDetail(w io.Writer, item memos.Memo) error {
	id := strings.TrimPrefix(item.Name, "memos/")
	content := strings.TrimSpace(item.Content)
	timeLine := formatMemoCreateTime(item.CreateTime)
	if timeLine != "" {
		_, err := fmt.Fprintf(w, "> %s\n%s\n%s\n", id, timeLine, content)
		return err
	}

	_, err := fmt.Fprintf(w, "> %s\n%s\n", id, content)
	return err
}

func formatMemoCreateTime(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}

	return parsed.Local().Format("2006-01-02 15:04")
}
