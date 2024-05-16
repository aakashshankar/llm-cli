package api

import (
	"fmt"
	"github.com/aakashshankar/claude-cli/anthropic"
	"os"
)

func PromptClaude(prompt string, stream bool) {
	anthropicConfig, err := anthropic.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	anthropicClient := anthropic.NewClient(anthropicConfig)
	err = anthropicClient.FetchCompletion(prompt, stream)
	if err != nil {
		fmt.Println("Error fetching completion:", err)
		os.Exit(1)
	}
}
