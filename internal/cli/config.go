package cli

import (
	"fmt"

	"github.com/rogeecn/memos-cli/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Inspect CLI configuration",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "Check active configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.LoadFromEnv()
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "MEMOS_URL: %s\nMEMOS_API_KEY: %s\nMEMOS_ADMIN_API_KEY: %s\nDEFAULT_TAG: %s\n",
				status(cfg.BaseURL), status(cfg.APIKey), status(cfg.AdminAPIKey), status(cfg.DefaultTag))
			return err
		},
	})

	return cmd
}

func status(value string) string {
	if value == "" {
		return "missing"
	}
	return "configured"
}
