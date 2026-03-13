package output

import (
	"fmt"
	"io"
)

func WriteUsers(w io.Writer, users []map[string]any) error {
	if len(users) == 0 {
		_, err := fmt.Fprintln(w, "No users found.")
		return err
	}
	for _, user := range users {
		name, _ := user["name"].(string)
		nickname, _ := user["nickname"].(string)
		if _, err := fmt.Fprintf(w, "%s\t%s\n", name, nickname); err != nil {
			return err
		}
	}
	return nil
}
