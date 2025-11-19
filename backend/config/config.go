package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig      `koanf:"app"`
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
	Auth     AuthConfig     `koanf:"auth"`
	AI       AIConfig       `koanf:"ai"`
	OAuth    OAuthConfig    `koanf:"oauth"`
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name        string `koanf:"name"`
	Version     string `koanf:"version"`
	Environment string `koanf:"environment"` // development, staging, production
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port            int      `koanf:"port"`
	Host            string   `koanf:"host"`
	ReadTimeout     int      `koanf:"read_timeout"`  // seconds
	WriteTimeout    int      `koanf:"write_timeout"` // seconds
	ShutdownTimeout int      `koanf:"shutdown_timeout"` // seconds
	AllowedOrigins  []string `koanf:"allowed_origins"` // CORS
}

// DatabaseConfig holds MongoDB configuration
type DatabaseConfig struct {
	URI      string `koanf:"uri"`
	Name     string `koanf:"name"`
	Timeout  int    `koanf:"timeout"`  // seconds
	PoolSize int    `koanf:"pool_size"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret           string `koanf:"jwt_secret"`
	JWTExpiration       int    `koanf:"jwt_expiration"`        // minutes
	RefreshExpiration   int    `koanf:"refresh_expiration"`    // days
	PasswordMinLength   int    `koanf:"password_min_length"`
}

// AIConfig holds AI service configuration
type AIConfig struct {
	Provider string `koanf:"provider"` // ollama, openai, localai
	BaseURL  string `koanf:"base_url"`
	APIKey   string `koanf:"api_key"`
	Model    string `koanf:"model"`
	Timeout  int    `koanf:"timeout"` // seconds
}

// OAuthConfig holds OAuth provider configuration
type OAuthConfig struct {
	Google GoogleOAuthConfig `koanf:"google"`
	GitHub GitHubOAuthConfig `koanf:"github"`
}

// GoogleOAuthConfig holds Google OAuth configuration
type GoogleOAuthConfig struct {
	ClientID     string `koanf:"client_id"`
	ClientSecret string `koanf:"client_secret"`
	RedirectURL  string `koanf:"redirect_url"`
}

// GitHubOAuthConfig holds GitHub OAuth configuration
type GitHubOAuthConfig struct {
	ClientID     string `koanf:"client_id"`
	ClientSecret string `koanf:"client_secret"`
	RedirectURL  string `koanf:"redirect_url"`
}

// Load loads configuration from config.yaml and environment variables
// Environment variables take precedence over config file values
func Load() (*Config, error) {
	k := koanf.New(".")

	// Load from config.yaml (if exists)
	if err := k.Load(file.Provider("config/config.yaml"), yaml.Parser()); err != nil {
		// Config file is optional, continue if it doesn't exist
		fmt.Printf("Warning: config.yaml not found, using defaults and env vars: %v\n", err)
	}

	// Load environment variables with prefix "HT_" (Hobby Tracker)
	// Environment variables override config file values
	// Example: HT_SERVER_PORT=8080 overrides server.port
	err := k.Load(env.Provider("HT_", ".", func(s string) string {
		// Convert HT_SERVER_PORT to server.port
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "HT_")), "_", ".", -1)
	}), nil)
	if err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}

	// Unmarshal into Config struct
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate configuration
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// validate checks if required configuration values are present
func validate(cfg *Config) error {
	if cfg.Server.Port == 0 {
		return fmt.Errorf("server.port is required")
	}
	if cfg.Database.URI == "" {
		return fmt.Errorf("database.uri is required")
	}
	if cfg.Database.Name == "" {
		return fmt.Errorf("database.name is required")
	}
	if cfg.Auth.JWTSecret == "" {
		return fmt.Errorf("auth.jwt_secret is required")
	}
	return nil
}
