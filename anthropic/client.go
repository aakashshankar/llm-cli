package anthropic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aakashshankar/claude-cli/session"
	"io"
	"net/http"
	"os"
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

func (c *Client) FetchCompletion(prompt string, stream bool, tokens int, model string, system string, clear bool) error {
	if clear {
		session.ClearSession()
	}
	s := session.NewSession()
	if err := s.LoadLatest(); err != nil {
		fmt.Println("Error loading session:", err)
		os.Exit(1)
	}
	req, err := marshalRequest(prompt, c, stream, tokens, model, system, s)
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

	if stream {
		response, ok := parseStreamingResponse(resp)
		if ok != nil {
			return ok
		}
		s.AddMessage("assistant", response)
	} else {
		response, ok := parseResponse(resp)
		if ok != nil {
			return ok
		}
		s.AddMessage("assistant", response.Content[0].Text)
	}
	err = s.Save()
	if err != nil {
		return err
	}
	return nil
}

func parseResponse(resp *http.Response) (*CompletionResponse, error) {
	var completion CompletionResponse
	err := json.NewDecoder(resp.Body).Decode(&completion)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response payload: %w", err)
	}
	response := completion.Content[0].Text
	fmt.Println(response)
	return &completion, nil
}

func parseStreamingResponse(resp *http.Response) (string, error) {
	reader := bufio.NewReader(resp.Body)
	var response string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("error reading response: %w", err)
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
						response += text
					}
				}
			}
		}
	}
	return response, nil
}

func marshalRequest(prompt string, c *Client, stream bool, tokens int, model string, system string,
	session *session.Session) (*http.Request, error) {
	messages := prependContext(prompt, session)
	payload := CompletionRequest{
		Model:     model,
		Messages:  messages,
		MaxTokens: tokens,
		Stream:    stream,
		System:    system,
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

func prependContext(prompt string, session *session.Session) []Message {
	var messages []Message
	for _, message := range session.Messages {
		messages = append(messages, Message{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	messages = append(messages, Message{
		Role:    "user",
		Content: prompt,
	})
	session.AddMessage("user", prompt)
	return messages
}
