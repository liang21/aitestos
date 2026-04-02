// Package config defines application configuration structures
package config

import (
	"fmt"
	"time"
)

// Config is the root configuration structure
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Milvus   MilvusConfig   `mapstructure:"milvus"`
	Storage  StorageConfig  `mapstructure:"storage"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// DatabaseConfig holds PostgreSQL database configuration
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LLMConfig holds LLM API configuration
type LLMConfig struct {
	Provider       string        `mapstructure:"provider"`
	APIKey         string        `mapstructure:"api_key"`
	BaseURL        string        `mapstructure:"base_url"`
	Model          string        `mapstructure:"model"`
	EmbeddingModel string        `mapstructure:"embedding_model"`
	Timeout        time.Duration `mapstructure:"timeout"`
	MaxRetries     int           `mapstructure:"max_retries"`
}

// MilvusConfig holds Milvus vector database configuration
type MilvusConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Database   string `mapstructure:"database"`
	Collection string `mapstructure:"collection"`
}

// StorageConfig holds object storage configuration
type StorageConfig struct {
	Provider  string `mapstructure:"provider"`
	Endpoint  string `mapstructure:"endpoint"`
	Region    string `mapstructure:"region"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

// RabbitMQConfig holds RabbitMQ configuration
type RabbitMQConfig struct {
	URL      string `mapstructure:"url"`
	Exchange string `mapstructure:"exchange"`
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	ExpireTime time.Duration `mapstructure:"expire_time"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string `mapstructure:"level"`
	JSON  bool   `mapstructure:"json"`
}

// Validate validates the entire configuration
func (c *Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config: %w", err)
	}
	if err := c.LLM.Validate(); err != nil {
		return fmt.Errorf("llm config: %w", err)
	}
	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("jwt config: %w", err)
	}
	return nil
}

// Validate validates ServerConfig
func (c *ServerConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf("shutdown timeout must be positive")
	}
	return nil
}

// Validate validates DatabaseConfig
func (c *DatabaseConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if c.Port <= 0 {
		return fmt.Errorf("port must be positive")
	}
	if c.User == "" {
		return fmt.Errorf("user cannot be empty")
	}
	if c.Database == "" {
		return fmt.Errorf("database name cannot be empty")
	}
	return nil
}

// ConnectionString returns PostgreSQL connection string
func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

// Validate validates LLMConfig
func (c *LLMConfig) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("api_key cannot be empty")
	}
	validProviders := map[string]bool{
		"deepseek": true,
		"openai":   true,
		"azure":    true,
	}
	if !validProviders[c.Provider] {
		return fmt.Errorf("invalid provider: %s", c.Provider)
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	return nil
}

// Validate validates JWTConfig
func (c *JWTConfig) Validate() error {
	if c.Secret == "" {
		return fmt.Errorf("secret cannot be empty")
	}
	if len(c.Secret) < 8 {
		return fmt.Errorf("secret must be at least 8 characters")
	}
	if c.ExpireTime < 5*time.Minute {
		return fmt.Errorf("expire time must be at least 5 minutes")
	}
	return nil
}
