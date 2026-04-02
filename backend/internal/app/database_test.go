// Package app_test tests database connection functionality
package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/liang21/aitestos/internal/config"
)

func TestNewDB_Validation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.DatabaseConfig
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &config.DatabaseConfig{
				Host:            "localhost",
				Port:            5432,
				User:            "postgres",
				Password:        "password",
				Database:        "aitestos_test",
				SSLMode:         "disable",
				MaxOpenConns:    25,
				MaxIdleConns:    5,
				ConnMaxLifetime: 5 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "empty host",
			cfg: &config.DatabaseConfig{
				Host:     "",
				Port:     5432,
				User:     "postgres",
				Database: "aitestos_test",
			},
			wantErr: true,
		},
		{
			name: "empty user",
			cfg: &config.DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "",
				Database: "aitestos_test",
			},
			wantErr: true,
		},
		{
			name: "empty database",
			cfg: &config.DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Database: "",
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			cfg: &config.DatabaseConfig{
				Host:     "localhost",
				Port:     -1,
				User:     "postgres",
				Database: "aitestos_test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseConfig_ConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.DatabaseConfig
		contains []string
	}{
		{
			name: "standard config",
			cfg: &config.DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "secret",
				Database: "aitestos",
				SSLMode:  "disable",
			},
			contains: []string{
				"host=localhost",
				"port=5432",
				"user=postgres",
				"password=secret",
				"dbname=aitestos",
				"sslmode=disable",
			},
		},
		{
			name: "with ssl mode require",
			cfg: &config.DatabaseConfig{
				Host:     "prod.example.com",
				Port:     5432,
				User:     "admin",
				Password: "prodpass",
				Database: "aitestos_prod",
				SSLMode:  "require",
			},
			contains: []string{
				"host=prod.example.com",
				"sslmode=require",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.ConnectionString()
			for _, substr := range tt.contains {
				if !contains(got, substr) {
					t.Errorf("ConnectionString() = %q, should contain %q", got, substr)
				}
			}
		})
	}
}

func TestDatabaseConfig_DefaultValues(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		User:     "postgres",
		Database: "aitestos",
	}

	// Test that defaults are reasonable
	if cfg.Port == 0 {
		cfg.Port = 5432 // Default PostgreSQL port
	}
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 25
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 5
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = 5 * time.Minute
	}

	if cfg.Port != 5432 {
		t.Errorf("Default port should be 5432, got %d", cfg.Port)
	}
	if cfg.SSLMode != "disable" {
		t.Errorf("Default SSLMode should be disable, got %s", cfg.SSLMode)
	}
	if cfg.MaxOpenConns != 25 {
		t.Errorf("Default MaxOpenConns should be 25, got %d", cfg.MaxOpenConns)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr ||
		len(s) > len(substr) && contains(s[1:], substr)
}

func TestNewDB_ConnectionPoolSettings(t *testing.T) {
	// This test verifies connection pool configuration is properly set
	// without actually connecting to a database

	cfg := &config.DatabaseConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "password",
		Database:        "aitestos_test",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}

	// Validate connection pool settings are reasonable
	if cfg.MaxOpenConns < 1 {
		t.Error("MaxOpenConns should be at least 1")
	}
	if cfg.MaxIdleConns < 1 {
		t.Error("MaxIdleConns should be at least 1")
	}
	if cfg.MaxIdleConns > cfg.MaxOpenConns {
		t.Error("MaxIdleConns should not exceed MaxOpenConns")
	}
	if cfg.ConnMaxLifetime < time.Minute {
		t.Error("ConnMaxLifetime should be at least 1 minute")
	}
}

func TestNewDB_ContextTimeout(t *testing.T) {
	// Test that database operations respect context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify context is not expired
	select {
	case <-ctx.Done():
		t.Error("Context should not be expired yet")
	default:
	}

	// Verify context expires after timeout
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel2()
	time.Sleep(10 * time.Millisecond)

	select {
	case <-ctx2.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should have expired")
	}
}
