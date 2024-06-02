package assist

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/aakashshankar/llm-cli/defaults"
	"github.com/aakashshankar/llm-cli/llm"
	"os"
)

//go:embed static/system_prompt.json
var systemPrompt []byte

type System struct {
	System string `json:"system"`
}

var system System

func init() {
	err := json.Unmarshal(systemPrompt, &system)
	if err != nil {
		fmt.Println("Error unmarshalling system prompt:", err)
		os.Exit(1)
	}
}

func Assist(prompt string) (string, error) {
	variant, ok := os.LookupEnv("DEFAULT_COMPLETER")
	if !ok {
		fmt.Println("DEFAULT_COMPLETER environment variable not set")
		os.Exit(1)
	}

	defaultModel, err := defaults.GetDefaultModel(variant)
	if err != nil {
		fmt.Println("Error getting default model:", err)
		os.Exit(1)
	}

	newLLM, err := llm.NewLLM(variant)
	if err != nil {
		return "", err
	}

	result, err := newLLM.Prompt(prompt, true, 1024, defaultModel, system.System, true)
	if err != nil {
		return "", err
	}
	return result, nil
}
