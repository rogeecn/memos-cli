package cli

import (
	"github.com/rogeecn/memos-cli/internal/config"
	"github.com/rogeecn/memos-cli/internal/memos"
)

func loadClientFromEnv() (*memos.Client, error) {
	cfg := config.LoadFromEnv()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return memos.NewClient(cfg.BaseURL, cfg.APIKey, cfg.AdminAPIKey), nil
}
