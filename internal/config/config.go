package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds the application configuration and API credentials
type Config struct {
	APIKey    string
	SecretKey string
}

// Load reads configuration from .env file (if present) and environment variables
func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		APIKey:    strings.TrimSpace(os.Getenv("BINANCE_API_KEY")),
		SecretKey: strings.TrimSpace(os.Getenv("BINANCE_SECRET_KEY")),
	}
}

// Validate ensures API credentials are provided and not default placeholder strings
func (c *Config) Validate() error {
	if c.APIKey == "" || c.SecretKey == "" {
		return fmt.Errorf("BINANCE_API_KEY or BINANCE_SECRET_KEY is missing. Please set them in your .env file or system environment variables")
	}

	if c.APIKey == "your_api_key_here" || c.SecretKey == "your_secret_key_here" {
		return fmt.Errorf("BINANCE_API_KEY or BINANCE_SECRET_KEY is set to placeholder values. Please update your .env file with real Binance API credentials")
	}

	return nil
}

// HasCredentials returns true if non-empty, non-placeholder credentials exist
func (c *Config) HasCredentials() bool {
	return c.Validate() == nil
}
