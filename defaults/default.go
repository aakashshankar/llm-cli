package defaults

import "fmt"

func GetDefaultModel(variant string) (string, error) {
	switch variant {
	case "claude":
		return "claude-3-sonnet-20240229", nil
	case "mistral":
		return "mistral-large-latest", nil
	default:
		return "", fmt.Errorf("unknown LLM type: %s", variant)
	}
}
