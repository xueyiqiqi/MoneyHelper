package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvFileIfPresentLoadsMissingValuesWithoutOverridingExistingEnv(t *testing.T) {
	originalBaseURL, hadBaseURL := os.LookupEnv("AI_BASE_URL")
	originalAPIKey, hadAPIKey := os.LookupEnv("AI_API_KEY")
	defer func() {
		if hadBaseURL {
			_ = os.Setenv("AI_BASE_URL", originalBaseURL)
		} else {
			_ = os.Unsetenv("AI_BASE_URL")
		}
		if hadAPIKey {
			_ = os.Setenv("AI_API_KEY", originalAPIKey)
		} else {
			_ = os.Unsetenv("AI_API_KEY")
		}
	}()

	if err := os.Setenv("AI_BASE_URL", "https://already-set.example.com"); err != nil {
		t.Fatalf("set AI_BASE_URL: %v", err)
	}
	if err := os.Unsetenv("AI_API_KEY"); err != nil {
		t.Fatalf("unset AI_API_KEY: %v", err)
	}

	envPath := filepath.Join(t.TempDir(), ".env")
	err := os.WriteFile(envPath, []byte("AI_BASE_URL=https://from-file.example.com\nAI_API_KEY=file-token\n"), 0o600)
	if err != nil {
		t.Fatalf("write env file: %v", err)
	}

	if err := LoadEnvFileIfPresent(envPath); err != nil {
		t.Fatalf("load env file: %v", err)
	}

	if got := os.Getenv("AI_BASE_URL"); got != "https://already-set.example.com" {
		t.Fatalf("expected existing env to win, got %q", got)
	}
	if got := os.Getenv("AI_API_KEY"); got != "file-token" {
		t.Fatalf("expected file token to be loaded, got %q", got)
	}
}
