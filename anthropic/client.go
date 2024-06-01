package anthropic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aakashshankar/llm-cli/session"
	"io"
	"net/http"
	"os"
	"strings"
)

const apiBaseURL = "https://api.anthropic.com"

type Client struct {
	config *Config
	client *http.Client
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		client: &http.Client{},
	}
}

func (c *Client) Prompt(prompt string, stream bool, tokens int, model string, system string, clear bool) (string, error) {
	if clear {
		session.ClearSession()
	}
	s := session.NewSession()
	if err := s.LoadLatest(); err != nil {
		fmt.Println("Error loading session:", err)
		os.Exit(1)
	}
	req, err := c.MarshalRequest(prompt, stream, tokens, model, system, s)
	if err != nil {
		return "", err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}
	var response string
	var ok error
	if stream {
		response, ok = c.ParseStreamingResponse(resp)
		if ok != nil {
			return "", ok
		}
		s.AddMessage("assistant", response)
	} else {
		completion, ok := c.ParseResponse(resp)
		if ok != nil {
			return "", ok
		}
		fmt.Println(completion)
		s.AddMessage("assistant", response)
	}
	err = s.Save()
	if err != nil {
		return "", err
	}
	return response, nil
}

func (c *Client) ParseResponse(resp *http.Response) (string, error) {
	var completion CompletionResponse
	err := json.NewDecoder(resp.Body).Decode(&completion)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response payload: %w", err)
	}
	return completion.Content[0].Text, nil
}

func (c *Client) ParseStreamingResponse(resp *http.Response) (string, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	reader := bufio.NewReader(resp.Body)
	var contentBuilder strings.Builder
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

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			data = strings.TrimSpace(data)

			var event CompletionStreamResponse
			err := json.Unmarshal([]byte(data), &event)
			if err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				continue
			}

			switch event.Type {
			case "content_block_delta":
				if event.Delta != nil {
					text := event.Delta.Text
					fmt.Print(text)
					contentBuilder.WriteString(text)
				}
			}
		}
	}
	return contentBuilder.String(), nil
}

func (c *Client) MarshalRequest(prompt string, stream bool, tokens int, model string, system string,
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
