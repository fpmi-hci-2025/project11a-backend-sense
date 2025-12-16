package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents application configuration
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	AI       AIConfig       `yaml:"ai"`
	Server   ServerConfig   `yaml:"server"`
	Media    MediaConfig    `yaml:"media"`
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// JWTConfig contains JWT token settings
type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expiry int    `yaml:"expiry"` // in seconds, default 86400 (1 day)
}

// AIConfig contains AI service settings
type AIConfig struct {
	ServiceURL string `yaml:"service_url"`
}

// ServerConfig contains server settings
type ServerConfig struct {
	Port int `yaml:"port"` // default 8080
}

// MediaConfig contains media upload settings
type MediaConfig struct {
	MaxFileSize int64 `yaml:"max_file_size"` // in bytes, default 10MB
}

// Load loads configuration from YAML file
func Load(configPath string) (*Config, error) {
	// #nosec G304 -- configPath is expected to be provided by the application, not user input
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.JWT.Expiry == 0 {
		config.JWT.Expiry = 86400 // 1 day
	}
	if config.Media.MaxFileSize == 0 {
		config.Media.MaxFileSize = 10 * 1024 * 1024 // 10MB
	}

	return &config, nil
}

// DSN returns PostgreSQL connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name)
}

