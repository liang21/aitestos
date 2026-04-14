// Package config provides configuration loading utilities
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Load reads configuration from file and environment variables
func Load(path string) (*Config, error) {
	v := viper.New()

	// Set config file
	v.SetConfigFile(path)

	// Allow environment variable substitution
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set defaults
	setDefaults(v)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	// Unmarshal to struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Validate
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

// setDefaults configures default values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.shutdown_timeout", "30s")

	// Database defaults
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", "5m")

	// Redis defaults
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	// LLM defaults
	v.SetDefault("llm.provider", "deepseek")
	v.SetDefault("llm.model", "deepseek-chat")
	v.SetDefault("llm.embedding_model", "deepseek-embedding")
	v.SetDefault("llm.timeout", "60s")
	v.SetDefault("llm.max_retries", 3)

	// Milvus defaults
	v.SetDefault("milvus.port", 19530)
	v.SetDefault("milvus.database", "aitestos")
	v.SetDefault("milvus.collection", "document_chunks")

	// Storage defaults
	v.SetDefault("storage.provider", "minio")
	v.SetDefault("storage.use_ssl", false)

	// RabbitMQ defaults
	v.SetDefault("rabbitmq.exchange", "aitestos.events")

	// JWT defaults
	v.SetDefault("jwt.expire_time", "2h")

	// Log defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.json", true)
}

// validate checks required configuration values
func validate(cfg *Config) error {
	if cfg.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if cfg.Database.User == "" {
		return fmt.Errorf("database.user is required")
	}
	if cfg.Database.Database == "" {
		return fmt.Errorf("database.database is required")
	}
	if cfg.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}
	if cfg.LLM.APIKey == "" {
		return fmt.Errorf("llm.api_key is required")
	}
	return nil
}
