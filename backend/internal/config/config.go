// Package config defines application configuration structures
package config

import (
	"time"
)

// Config is the root configuration structure
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	LLM       LLMConfig       `mapstructure:"llm"`
	Milvus    MilvusConfig    `mapstructure:"milvus"`
	Storage   StorageConfig   `mapstructure:"storage"`
	RabbitMQ  RabbitMQConfig  `mapstructure:"rabbitmq"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Log       LogConfig       `mapstructure:"log"`
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
