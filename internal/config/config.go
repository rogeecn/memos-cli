package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	BaseURL     string
	APIKey      string
	AdminAPIKey string
	DefaultTag  string
}

func LoadFromEnv() Config {
	return Config{
		BaseURL:     strings.TrimSpace(os.Getenv("MEMOS_URL")),
		APIKey:      strings.TrimSpace(os.Getenv("MEMOS_API_KEY")),
		AdminAPIKey: strings.TrimSpace(os.Getenv("MEMOS_ADMIN_API_KEY")),
		DefaultTag:  strings.TrimSpace(os.Getenv("DEFAULT_TAG")),
	}
}

func (c Config) Validate() error {
	if c.BaseURL == "" {
		return errors.New("MEMOS_URL is required")
	}
	if c.APIKey == "" {
		return errors.New("MEMOS_API_KEY is required")
	}
	return nil
}

func (c Config) ValidateAdmin() error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.AdminAPIKey == "" {
		return errors.New("MEMOS_ADMIN_API_KEY is required")
	}
	return nil
}
