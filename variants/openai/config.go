package openai

import (
	"fmt"
	"os"
)

type Config struct {
	APIKey string
}

func LoadConfig() (*Config, error) {
	apiKey, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	return &Config{
		APIKey: apiKey,
	}, nil
}
