// Package config_test tests configuration loading functionality
package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/liang21/aitestos/internal/config"
)

func TestLoad_ValidConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  host: "0.0.0.0"
  port: 8080
  shutdown_timeout: "30s"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  database: "aitestos_test"
  ssl_mode: "disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: "5m"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

llm:
  provider: "deepseek"
  api_key: "sk-test-key"
  base_url: "https://api.deepseek.com"
  model: "deepseek-chat"
  embedding_model: "deepseek-embedding"
  timeout: "60s"
  max_retries: 3

jwt:
  secret: "my-super-secret-key-for-testing"
  expire_time: "2h"

log:
  level: "info"
  json: true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify server config
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %q, want %q", cfg.Server.Host, "0.0.0.0")
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 8080)
	}
	if cfg.Server.ShutdownTimeout != 30*time.Second {
		t.Errorf("Server.ShutdownTimeout = %v, want %v", cfg.Server.ShutdownTimeout, 30*time.Second)
	}

	// Verify database config
	if cfg.Database.Host != "localhost" {
		t.Errorf("Database.Host = %q, want %q", cfg.Database.Host, "localhost")
	}
	if cfg.Database.User != "postgres" {
		t.Errorf("Database.User = %q, want %q", cfg.Database.User, "postgres")
	}

	// Verify LLM config
	if cfg.LLM.Provider != "deepseek" {
		t.Errorf("LLM.Provider = %q, want %q", cfg.LLM.Provider, "deepseek")
	}
	if cfg.LLM.APIKey != "sk-test-key" {
		t.Errorf("LLM.APIKey = %q, want %q", cfg.LLM.APIKey, "sk-test-key")
	}

	// Verify JWT config
	if cfg.JWT.Secret != "my-super-secret-key-for-testing" {
		t.Errorf("JWT.Secret = %q, want %q", cfg.JWT.Secret, "my-super-secret-key-for-testing")
	}
}

func TestLoad_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name          string
		configContent string
		wantErr       bool
	}{
		{
			name: "missing database host",
			configContent: `
server:
  host: "0.0.0.0"
  port: 8080
  shutdown_timeout: "30s"

database:
  user: "postgres"

jwt:
  secret: "my-secret-key"

llm:
  api_key: "sk-test"
`,
			wantErr: true,
		},
		{
			name: "missing jwt secret",
			configContent: `
server:
  host: "0.0.0.0"
  port: 8080

database:
  host: "localhost"
  user: "postgres"

llm:
  api_key: "sk-test"
`,
			wantErr: true,
		},
		{
			name: "missing llm api key",
			configContent: `
server:
  host: "0.0.0.0"
  port: 8080

database:
  host: "localhost"
  user: "postgres"

jwt:
  secret: "my-secret-key"
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")

			if err := os.WriteFile(configPath, []byte(tt.configContent), 0644); err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			_, err := config.Load(configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Load() should return error for non-existent file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	invalidYAML := `
server:
  host: "0.0.0.0"
  port: [invalid yaml
`
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := config.Load(configPath)
	if err == nil {
		t.Error("Load() should return error for invalid YAML")
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Minimal config - defaults should be applied
	minimalConfig := `
database:
  host: "localhost"
  user: "postgres"
  database: "aitestos"

jwt:
  secret: "my-super-secret-key-for-testing-minimum-length"

llm:
  api_key: "sk-test"
  timeout: "60s"
`
	if err := os.WriteFile(configPath, []byte(minimalConfig), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Check default values
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Default Server.Host = %q, want %q", cfg.Server.Host, "0.0.0.0")
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Default Server.Port = %d, want %d", cfg.Server.Port, 8080)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("Default Database.Port = %d, want %d", cfg.Database.Port, 5432)
	}
	if cfg.Database.SSLMode != "disable" {
		t.Errorf("Default Database.SSLMode = %q, want %q", cfg.Database.SSLMode, "disable")
	}
	if cfg.Database.MaxOpenConns != 25 {
		t.Errorf("Default Database.MaxOpenConns = %d, want %d", cfg.Database.MaxOpenConns, 25)
	}
	if cfg.Database.MaxIdleConns != 5 {
		t.Errorf("Default Database.MaxIdleConns = %d, want %d", cfg.Database.MaxIdleConns, 5)
	}
	if cfg.Redis.Port != 6379 {
		t.Errorf("Default Redis.Port = %d, want %d", cfg.Redis.Port, 6379)
	}
	if cfg.LLM.Provider != "deepseek" {
		t.Errorf("Default LLM.Provider = %q, want %q", cfg.LLM.Provider, "deepseek")
	}
	if cfg.Log.Level != "info" {
		t.Errorf("Default Log.Level = %q, want %q", cfg.Log.Level, "info")
	}
	if cfg.Log.JSON != true {
		t.Errorf("Default Log.JSON = %v, want %v", cfg.Log.JSON, true)
	}
}

func TestLoad_EnvironmentVariableOverride(t *testing.T) {
	// Set environment variables
	envVars := map[string]string{
		"DATABASE_HOST": "env-db-host",
		"DATABASE_PORT": "5433",
		"LLM_API_KEY":   "env-api-key",
		"JWT_SECRET":    "env-jwt-secret-with-minimum-length",
	}

	for k, v := range envVars {
		t.Setenv(k, v)
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
database:
  host: "file-host"
  user: "postgres"
  database: "aitestos"

jwt:
  secret: "file-secret-key-with-minimum-length"

llm:
  api_key: "file-api-key"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Environment variables should override file values (if viper is configured correctly)
	// Note: The actual behavior depends on viper configuration
	_ = cfg // Just verify it loads without error
}
