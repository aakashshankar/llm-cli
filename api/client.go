package api

import (
	"bufio"
	"fmt"
	"github.com/aakashshankar/llm-cli/anthropic"
	"github.com/aakashshankar/llm-cli/session"
	"os"
)

func Prompt(llm string, prompt string, stream bool, tokens int, model string, system string, clear bool) {
	switch llm {
	case "claude":
		PromptClaude(prompt, stream, tokens, model, system, clear)
	default:
		fmt.Println("Unknown LLM:", llm)
		os.Exit(1)
	}
}

func Chat(llm string) {
	reader := bufio.NewReader(os.Stdin)
	session.ClearSession()
	switch llm {
	case "claude":
		for {
			fmt.Print("> ")
			text, _ := reader.ReadString('\n')
			if text == "exit\n" {
				break
			}
			PromptClaude(text, true, 1024, "claude-3-sonnet-20240229", "You are a helpful assistant.", false)
			//fmt.Println()
		}
	default:
		fmt.Println("Unknown LLM:", llm)
		os.Exit(1)
	}
}

func PromptClaude(prompt string, stream bool, tokens int, model string, system string, clear bool) {

	anthropicConfig, err := anthropic.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	anthropicClient := anthropic.NewClient(anthropicConfig)
	_, err = anthropicClient.Prompt(prompt, stream, tokens, model, system, clear)
	if err != nil {
		fmt.Println("Error fetching completion:", err)
		os.Exit(1)
	}
}
