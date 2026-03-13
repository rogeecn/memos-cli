package cli

import (
	"github.com/rogeecn/memos-cli/internal/memos"
	"github.com/rogeecn/memos-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCommentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Manage memo comments",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "create <memo-id> <content>",
		Short: "Create a comment for a memo",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := loadClientFromEnv()
			if err != nil {
				return err
			}
			memo, err := client.CreateComment(args[0], memos.MemoPayload{Content: args[1], Visibility: "PRIVATE"})
			if err != nil {
				return err
			}
			return output.WriteMemoDetail(cmd.OutOrStdout(), memo)
		},
	})

	return cmd
}
