package api

import (
	"bufio"
	"fmt"
	"github.com/aakashshankar/llm-cli/llm"
	"github.com/aakashshankar/llm-cli/session"
	"os"
)

func Prompt(llmType string, prompt string, stream bool, tokens int, model string, system string, clear bool) {
	newLLM, err := llm.NewLLM(llmType)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	_, err = newLLM.Prompt(prompt, stream, tokens, model, system, clear)
	if err != nil {
		fmt.Println("Error fetching completion:", err)
		os.Exit(1)
	}
}

func Chat(llmType string) {
	newLLM, err := llm.NewLLM(llmType)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	session.ClearSession()
	var model string
	switch llmType {
	case "claude":
		model = "claude-3-sonnet-20240229"
	case "mistral":
		model = "mistral-large-latest"
	default:
		fmt.Println("Unknown LLM:", llmType)
		os.Exit(1)
	}

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		if text == "exit\n" {
			break
		}
		_, err := newLLM.Prompt(text, true, 1024, model, "You are a helpful assistant.", false)
		if err != nil {
			fmt.Println("Error fetching completion:", err)
			continue
		}
		fmt.Println()
	}
}
