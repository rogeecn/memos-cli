package cli

import (
	"fmt"

	"github.com/rogeecn/memos-cli/internal/memos"
	"github.com/rogeecn/memos-cli/internal/output"
	"github.com/spf13/cobra"
)

func newMemoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memo",
		Short: "Read and manage memos",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List memos",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := loadClientFromEnv()
			if err != nil {
				return err
			}
			response, err := client.ListMemos(memos.ListMemosParams{})
			if err != nil {
				return err
			}
			if getCommandContext(cmd.Context()).jsonOutput {
				return output.WriteJSON(cmd.OutOrStdout(), response)
			}
			return output.WriteMemoList(cmd.OutOrStdout(), response.Memos)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get <memo-id>",
		Short: "Get a memo by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := loadClientFromEnv()
			if err != nil {
				return err
			}
			memo, err := client.GetMemo(args[0])
			if err != nil {
				return err
			}
			if getCommandContext(cmd.Context()).jsonOutput {
				return output.WriteJSON(cmd.OutOrStdout(), memo)
			}
			return output.WriteMemoDetail(cmd.OutOrStdout(), memo)
		},
	})

	cmd.AddCommand(newMemoCreateCommand())
	cmd.AddCommand(newMemoUpdateCommand())
	cmd.AddCommand(newMemoDeleteCommand())

	return cmd
}

func newSearchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search memos by content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := loadClientFromEnv()
			if err != nil {
				return err
			}
			filter := fmt.Sprintf("content.contains('%s')", args[0])
			response, err := client.ListMemos(memos.ListMemosParams{Filter: filter})
			if err != nil {
				return err
			}
			if getCommandContext(cmd.Context()).jsonOutput {
				return output.WriteJSON(cmd.OutOrStdout(), response)
			}
			return output.WriteMemoList(cmd.OutOrStdout(), response.Memos)
		},
	}
}

func newFilterCommand() *cobra.Command {
	var expr string
	cmd := &cobra.Command{
		Use:   "filter",
		Short: "Filter memos with a CEL expression",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := loadClientFromEnv()
			if err != nil {
				return err
			}
			response, err := client.ListMemos(memos.ListMemosParams{Filter: expr})
			if err != nil {
				return err
			}
			if getCommandContext(cmd.Context()).jsonOutput {
				return output.WriteJSON(cmd.OutOrStdout(), response)
			}
			return output.WriteMemoList(cmd.OutOrStdout(), response.Memos)
		},
	}
	cmd.Flags().StringVar(&expr, "expr", "", "CEL filter expression")
	_ = cmd.MarkFlagRequired("expr")
	return cmd
}
