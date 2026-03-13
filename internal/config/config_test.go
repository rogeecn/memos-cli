package config

import "testing"

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
