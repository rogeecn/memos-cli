package cli

import (
	"github.com/rogeecn/memos-cli/internal/input"
	"github.com/rogeecn/memos-cli/internal/memos"
	"github.com/rogeecn/memos-cli/internal/output"
	"github.com/spf13/cobra"
)

func newTagCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage memo tags",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "remove <memo-id> <tag>",
		Short: "Remove a tag from a memo",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := loadClientFromEnv()
			if err != nil {
				return err
			}
			memo, err := client.GetMemo(args[0])
			if err != nil {
				return err
			}
			updated, err := client.UpdateMemo(args[0], memos.UpdateMemoPayload{Content: input.RemoveTag(memo.Content, args[1])})
			if err != nil {
				return err
			}
			return output.WriteMemoDetail(cmd.OutOrStdout(), updated)
		},
	})

	return cmd
}
