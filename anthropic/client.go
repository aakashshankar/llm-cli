package anthropic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func (c *Client) FetchCompletion(prompt string, stream bool) error {
	req, err := marshalRequest(prompt, c, stream)
	if err != nil {
		return err
	}
	// Send the request and handle the response.
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	var ok error
	if stream {
		ok = parseStreamingResponse(resp)
		if ok != nil {
			return ok
		}
	} else {
		ok = parseResponse(resp)
		if ok != nil {
			return ok
		}
	}

	return nil
}

func parseResponse(resp *http.Response) error {
	var completion CompletionResponse
	err := json.NewDecoder(resp.Body).Decode(&completion)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response payload: %w", err)
	}
	fmt.Println(completion.Content[0].Text)
	return nil
}

func parseStreamingResponse(resp *http.Response) error {
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading response: %w", err)
		}
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Check if the line starts with "data: "
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			data = strings.TrimSpace(data)

			// Unmarshal the JSON data
			var eventData map[string]interface{}
			err := json.Unmarshal([]byte(data), &eventData)
			if err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				continue
			}

			if eventType, ok := eventData["type"].(string); ok && eventType == "content_block_delta" {
				if delta, ok := eventData["delta"].(map[string]interface{}); ok {
					if text, ok := delta["text"].(string); ok {
						fmt.Print(text)
					}
				}
			}
		}
	}
	return nil
}

func marshalRequest(prompt string, c *Client, stream bool) (*http.Request, error) {
	payload := map[string]interface{}{
		"max_tokens": 1024,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"stream": stream,
		"model":  "claude-3-sonnet-20240229",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/messages", apiBaseURL), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.config.APIKey)
	req.Header.Set("anthropic-version", c.config.AnthropicVersion)
	return req, nil
}
