package gemini

import (
	"fmt"
	"os"
)

type Config struct {
	APIKey string
}

func LoadConfig() (*Config, error) {
	apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	return &Config{
		APIKey: apiKey,
	}, nil
}
