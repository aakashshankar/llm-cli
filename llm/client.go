package llm

type Client interface {
	Prompt(prompt string, stream bool, tokens int, model string, system string, clear bool) (string, error)
}
