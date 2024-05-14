package api

import (
	"fmt"
	"github.com/aakashshankar/claude-cli/anthropic"
	"os"
)

func PromptClaude(prompt string) string {
	anthropicConfig, err := anthropic.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	anthropicClient := anthropic.NewClient(anthropicConfig)
	response, err := anthropicClient.FetchCompletion(prompt)
	if err != nil {
		fmt.Println("Error fetching completion:", err)
		os.Exit(1)
	}
	return response.Content[0].Text
}
