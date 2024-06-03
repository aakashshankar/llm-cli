package openai

type CompletionRequest struct {
	Model            string      `json:"model"`
	Messages         []Message   `json:"messages"`
	Temperature      float64     `json:"temperature,omitempty"`
	TopP             float64     `json:"top_p,omitempty"`
	FrequencyPenalty float64     `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64     `json:"presence_penalty,omitempty"`
	LogitBias        interface{} `json:"logit_bias,omitempty"`
	N                int         `json:"n,omitempty"`
	LogProbs         bool        `json:"logprobs,omitempty"`
	TopLogProbs      int         `json:"top_logprobs,omitempty"`
	ResponseFormat   string      `json:"response_format,omitempty"`
	MaxTokens        int         `json:"max_tokens,omitempty"`
	Stream           bool        `json:"stream,omitempty"`
	StreamOptions    interface{} `json:"stream_options,omitempty"`
	Seed             int         `json:"seed,omitempty"`
	Stop             []string    `json:"stop,omitempty"`
	Tools            []struct {
		Type     string `json:"type"`
		Function string `json:"function"`
	} `json:"tools,omitempty"`
	User string `json:"user,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	}
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

type CompletionStreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason interface{} `json:"finish_reason"`
		Logprobs     interface{} `json:"logprobs"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
