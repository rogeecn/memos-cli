package config

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL     string
	APIKey      string
	AdminAPIKey string
	DefaultTag  string
}

func LoadFromEnv() Config {
	dotenvValues, err := godotenv.Read()
	if err != nil {
		dotenvValues = map[string]string{}
	}

	return Config{
		BaseURL:     loadValue("MEMOS_URL", dotenvValues),
		APIKey:      loadValue("MEMOS_API_KEY", dotenvValues),
		AdminAPIKey: loadValue("MEMOS_ADMIN_API_KEY", dotenvValues),
		DefaultTag:  loadValue("DEFAULT_TAG", dotenvValues),
	}
}

func loadValue(key string, dotenvValues map[string]string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}

	return strings.TrimSpace(dotenvValues[key])
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
