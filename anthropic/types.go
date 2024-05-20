package anthropic

// TextContent represents the text content in the response.
type TextContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

// Usage represents the token usage in the response.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// CompletionResponse represents the response from the Anthropic API for a completion request.
type CompletionResponse struct {
	Content      []TextContent `json:"content"`
	ID           string        `json:"id"`
	Model        string        `json:"model"`
	Role         string        `json:"role"`
	StopReason   string        `json:"stop_reason"`
	StopSequence interface{}   `json:"stop_sequence"` // Can be null
	Type         string        `json:"type"`
	Usage        Usage         `json:"usage"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Metadata struct {
	UserId string `json:"user_id"`
}

type CompletionRequest struct {
	Model         string    `json:"model"`
	Messages      []Message `json:"messages"`
	MaxTokens     int       `json:"max_tokens"`
	Metadata      Metadata  `json:"metadata"`
	Stream        bool      `json:"stream"`
	StopSequences []string  `json:"stop_sequences"`
	System        string    `json:"system"`
	Temperature   float64   `json:"temperature"`
	TopK          int       `json:"top_k"`
	TopP          float64   `json:"top_p"`
}
