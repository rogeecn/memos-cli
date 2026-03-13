package cli

import (
	"github.com/rogeecn/memos-cli/internal/output"
	"github.com/spf13/cobra"
)

func newUserCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Inspect users",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List users via admin API",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := loadClientFromEnv()
			if err != nil {
				return err
			}
			response, err := client.ListUsers()
			if err != nil {
				return err
			}
			if getCommandContext(cmd.Context()).jsonOutput {
				return output.WriteJSON(cmd.OutOrStdout(), response)
			}
			return output.WriteUsers(cmd.OutOrStdout(), response.Users)
		},
	})

	return cmd
}
