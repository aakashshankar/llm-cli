package anthropic

import (
	"fmt"
	"os"
)

// Config represents the configuration for the Anthropic API client.
type Config struct {
	APIKey           string
	AnthropicVersion string
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig() (*Config, error) {
	apiKey, ok := os.LookupEnv("ANTHROPIC_API_KEY")
	if !ok {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}
	version := "2023-06-01"
	return &Config{
		APIKey:           apiKey,
		AnthropicVersion: version,
	}, nil
}
