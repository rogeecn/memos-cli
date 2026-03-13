package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromEnvReadsExpectedVariables(t *testing.T) {
	t.Setenv("MEMOS_URL", "https://memos.example.com")
	t.Setenv("MEMOS_API_KEY", "token-123")
	t.Setenv("MEMOS_ADMIN_API_KEY", "admin-456")
	t.Setenv("DEFAULT_TAG", "cli")

	cfg := LoadFromEnv()

	if cfg.BaseURL != "https://memos.example.com" {
		t.Fatalf("expected base url from env, got %q", cfg.BaseURL)
	}
	if cfg.APIKey != "token-123" {
		t.Fatalf("expected api key from env, got %q", cfg.APIKey)
	}
	if cfg.AdminAPIKey != "admin-456" {
		t.Fatalf("expected admin api key from env, got %q", cfg.AdminAPIKey)
	}
	if cfg.DefaultTag != "cli" {
		t.Fatalf("expected default tag from env, got %q", cfg.DefaultTag)
	}
}

func TestLoadFromEnvReadsDotenvWhenEnvironmentMissing(t *testing.T) {
	withWorkingDirectory(t, func(dir string) {
		writeDotenvFile(t, dir, "MEMOS_URL=http://from-dotenv\nMEMOS_API_KEY=dotenv-key\nMEMOS_ADMIN_API_KEY=dotenv-admin\nDEFAULT_TAG=dotenv-tag\n")

		cfg := LoadFromEnv()

		if cfg.BaseURL != "http://from-dotenv" {
			t.Fatalf("expected base url from .env, got %q", cfg.BaseURL)
		}
		if cfg.APIKey != "dotenv-key" {
			t.Fatalf("expected api key from .env, got %q", cfg.APIKey)
		}
		if cfg.AdminAPIKey != "dotenv-admin" {
			t.Fatalf("expected admin api key from .env, got %q", cfg.AdminAPIKey)
		}
		if cfg.DefaultTag != "dotenv-tag" {
			t.Fatalf("expected default tag from .env, got %q", cfg.DefaultTag)
		}
	})
}

func TestLoadFromEnvPrefersSystemEnvironmentOverDotenv(t *testing.T) {
	t.Setenv("MEMOS_URL", "https://from-env")
	t.Setenv("MEMOS_API_KEY", "env-key")
	t.Setenv("MEMOS_ADMIN_API_KEY", "env-admin")
	t.Setenv("DEFAULT_TAG", "env-tag")

	withWorkingDirectory(t, func(dir string) {
		writeDotenvFile(t, dir, "MEMOS_URL=http://from-dotenv\nMEMOS_API_KEY=dotenv-key\nMEMOS_ADMIN_API_KEY=dotenv-admin\nDEFAULT_TAG=dotenv-tag\n")

		cfg := LoadFromEnv()

		if cfg.BaseURL != "https://from-env" {
			t.Fatalf("expected base url from env, got %q", cfg.BaseURL)
		}
		if cfg.APIKey != "env-key" {
			t.Fatalf("expected api key from env, got %q", cfg.APIKey)
		}
		if cfg.AdminAPIKey != "env-admin" {
			t.Fatalf("expected admin api key from env, got %q", cfg.AdminAPIKey)
		}
		if cfg.DefaultTag != "env-tag" {
			t.Fatalf("expected default tag from env, got %q", cfg.DefaultTag)
		}
	})
}

func TestValidateRequiresBaseURLAndAPIKey(t *testing.T) {
	cfg := Config{}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestValidateAdminRequiresAdminKey(t *testing.T) {
	cfg := Config{BaseURL: "https://memos.example.com", APIKey: "token-123"}

	err := cfg.ValidateAdmin()
	if err == nil {
		t.Fatal("expected admin validation error, got nil")
	}
}

func withWorkingDirectory(t *testing.T, run func(dir string)) {
	t.Helper()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	run(tempDir)
}

func writeDotenvFile(t *testing.T, dir string, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0o600); err != nil {
		t.Fatalf("write .env file: %v", err)
	}
}
