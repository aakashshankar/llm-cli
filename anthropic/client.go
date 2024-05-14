package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const apiBaseURL = "https://api.anthropic.com"

// Client represents the Anthropic API client.
type Client struct {
	config *Config
	client *http.Client
}

// NewClient creates a new Anthropic API client with the provided configuration.
func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		client: &http.Client{},
	}
}

// FetchCompletion fetches a completion from the Anthropic API using Claude.
func (c *Client) FetchCompletion(prompt string) (*CompletionResponse, error) {
	// Prepare the request payload.
	payload := map[string]interface{}{
		"max_tokens": 1024,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"model": "claude-3-sonnet-20240229",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Create the request.
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/messages", apiBaseURL), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.config.APIKey)
	req.Header.Set("anthropic-version", c.config.AnthropicVersion)

	// Send the request and handle the response.
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, body)
	}

	var completionResponse CompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&completionResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response payload: %w", err)
	}

	return &completionResponse, nil
}
