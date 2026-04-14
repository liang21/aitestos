// Package config_test tests configuration structures
package config_test

import (
	"testing"
	"time"

	"github.com/liang21/aitestos/internal/config"
)

func TestConfig_DefaultValues(t *testing.T) {
	tests := []struct {
		name     string
		config   config.Config
		field    string
		expected any
	}{
		{
			name:     "server host default",
			config:   config.Config{Server: config.ServerConfig{Host: "0.0.0.0"}},
			field:    "Server.Host",
			expected: "0.0.0.0",
		},
		{
			name:     "server port default",
			config:   config.Config{Server: config.ServerConfig{Port: 8080}},
			field:    "Server.Port",
			expected: 8080,
		},
		{
			name:     "shutdown timeout default",
			config:   config.Config{Server: config.ServerConfig{ShutdownTimeout: 30 * time.Second}},
			field:    "Server.ShutdownTimeout",
			expected: 30 * time.Second,
		},
		{
			name:     "database max open conns",
			config:   config.Config{Database: config.DatabaseConfig{MaxOpenConns: 25}},
			field:    "Database.MaxOpenConns",
			expected: 25,
		},
		{
			name:     "database max idle conns",
			config:   config.Config{Database: config.DatabaseConfig{MaxIdleConns: 5}},
			field:    "Database.MaxIdleConns",
			expected: 5,
		},
		{
			name:     "jwt expire time",
			config:   config.Config{JWT: config.JWTConfig{ExpireTime: 2 * time.Hour}},
			field:    "JWT.ExpireTime",
			expected: 2 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.field {
			case "Server.Host":
				if tt.config.Server.Host != tt.expected.(string) {
					t.Errorf("Server.Host = %v, want %v", tt.config.Server.Host, tt.expected)
				}
			case "Server.Port":
				if tt.config.Server.Port != tt.expected.(int) {
					t.Errorf("Server.Port = %v, want %v", tt.config.Server.Port, tt.expected)
				}
			case "Server.ShutdownTimeout":
				if tt.config.Server.ShutdownTimeout != tt.expected.(time.Duration) {
					t.Errorf("Server.ShutdownTimeout = %v, want %v", tt.config.Server.ShutdownTimeout, tt.expected)
				}
			case "Database.MaxOpenConns":
				if tt.config.Database.MaxOpenConns != tt.expected.(int) {
					t.Errorf("Database.MaxOpenConns = %v, want %v", tt.config.Database.MaxOpenConns, tt.expected)
				}
			case "Database.MaxIdleConns":
				if tt.config.Database.MaxIdleConns != tt.expected.(int) {
					t.Errorf("Database.MaxIdleConns = %v, want %v", tt.config.Database.MaxIdleConns, tt.expected)
				}
			case "JWT.ExpireTime":
				if tt.config.JWT.ExpireTime != tt.expected.(time.Duration) {
					t.Errorf("JWT.ExpireTime = %v, want %v", tt.config.JWT.ExpireTime, tt.expected)
				}
			}
		})
	}
}

func TestServerConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  config.ServerConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: config.ServerConfig{
				Host:            "0.0.0.0",
				Port:            8080,
				ShutdownTimeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "invalid port - too low",
			config: config.ServerConfig{
				Host:            "0.0.0.0",
				Port:            0,
				ShutdownTimeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid port - too high",
			config: config.ServerConfig{
				Host:            "0.0.0.0",
				Port:            70000,
				ShutdownTimeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "empty host",
			config: config.ServerConfig{
				Host:            "",
				Port:            8080,
				ShutdownTimeout: 30 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseConfig_ConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		config   config.DatabaseConfig
		expected string
	}{
		{
			name: "standard connection string",
			config: config.DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "secret",
				Database: "aitestos",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=postgres password=secret dbname=aitestos sslmode=disable",
		},
		{
			name: "with ssl mode",
			config: config.DatabaseConfig{
				Host:     "prod-db.example.com",
				Port:     5432,
				User:     "admin",
				Password: "prodpass",
				Database: "aitestos_prod",
				SSLMode:  "require",
			},
			expected: "host=prod-db.example.com port=5432 user=admin password=prodpass dbname=aitestos_prod sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.ConnectionString()
			if got != tt.expected {
				t.Errorf("DatabaseConfig.ConnectionString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLLMConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  config.LLMConfig
		wantErr bool
	}{
		{
			name: "valid deepseek config",
			config: config.LLMConfig{
				Provider:       "deepseek",
				APIKey:         "sk-xxx",
				BaseURL:        "https://api.deepseek.com",
				Model:          "deepseek-chat",
				EmbeddingModel: "deepseek-embedding",
				Timeout:        60 * time.Second,
				MaxRetries:     3,
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			config: config.LLMConfig{
				Provider:       "deepseek",
				APIKey:         "",
				BaseURL:        "https://api.deepseek.com",
				Model:          "deepseek-chat",
				EmbeddingModel: "deepseek-embedding",
				Timeout:        60 * time.Second,
				MaxRetries:     3,
			},
			wantErr: true,
		},
		{
			name: "invalid provider",
			config: config.LLMConfig{
				Provider:       "invalid",
				APIKey:         "sk-xxx",
				BaseURL:        "https://api.example.com",
				Model:          "model",
				EmbeddingModel: "embedding",
				Timeout:        60 * time.Second,
				MaxRetries:     3,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("LLMConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJWTConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  config.JWTConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: config.JWTConfig{
				Secret:     "my-secret-key",
				ExpireTime: 2 * time.Hour,
			},
			wantErr: false,
		},
		{
			name: "empty secret",
			config: config.JWTConfig{
				Secret:     "",
				ExpireTime: 2 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "secret too short",
			config: config.JWTConfig{
				Secret:     "short",
				ExpireTime: 2 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "expire time too short",
			config: config.JWTConfig{
				Secret:     "my-secret-key",
				ExpireTime: 1 * time.Minute,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("JWTConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
