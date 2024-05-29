package llm

import (
	"fmt"
	"github.com/aakashshankar/llm-cli/anthropic"
	"github.com/aakashshankar/llm-cli/mistral"
	"github.com/aakashshankar/llm-cli/session"
	"net/http"
)

type LLM interface {
	Prompt(prompt string, stream bool, tokens int, model string, system string, clear bool) (string, error)
	MarshalRequest(prompt string, stream bool, tokens int, model string, system string, session *session.Session) (*http.Request, error)
	ParseResponse(resp *http.Response) (string, error)
	ParseStreamingResponse(resp *http.Response) (string, error)
}

func NewLLM(llmType string) (LLM, error) {
	switch llmType {
	case "claude":
		config, err := anthropic.LoadConfig()
		if err != nil {
			return nil, fmt.Errorf("error loading config for Claude: %w", err)
		}
		return anthropic.NewClient(config), nil
	case "mistral":
		config, err := mistral.LoadConfig()
		if err != nil {
			return nil, fmt.Errorf("error loading config for Mistral: %w", err)
		}
		return mistral.NewClient(config), nil
	default:
		return nil, fmt.Errorf("unknown LLM type: %s", llmType)
	}
}
