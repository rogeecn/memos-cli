package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	ctx := &commandContext{}
	cmd := &cobra.Command{
		Use:   "memos",
		Short: "Memos CLI",
		Long:  "Memos CLI provides direct terminal access to common Memos operations.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SetContext(withCommandContext(cmd.Context(), ctx))
		},
	}

	cmd.PersistentFlags().BoolVar(&ctx.jsonOutput, "json", false, "Output raw JSON")
	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newMemoCommand())
	cmd.AddCommand(newSearchCommand())
	cmd.AddCommand(newFilterCommand())
	cmd.AddCommand(newCommentCommand())
	cmd.AddCommand(newTagCommand())
	cmd.AddCommand(newUserCommand())

	return cmd
}
