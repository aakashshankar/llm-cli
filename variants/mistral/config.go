package mistral

import (
	"fmt"
	"os"
)

type Config struct {
	APIKey string
}

func LoadConfig() (*Config, error) {
	apiKey, ok := os.LookupEnv("MISTRAL_API_KEY")
	if !ok {
		return nil, fmt.Errorf("MISTRAL_API_KEY environment variable not set")
	}

	return &Config{
		APIKey: apiKey,
	}, nil
}
