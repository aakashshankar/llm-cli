package mistral

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aakashshankar/llm-cli/highlight"
	"github.com/aakashshankar/llm-cli/session"
	"io"
	"net/http"
	"os"
	"strings"
)

const v1apiBaseURL = "https://api.mistral.ai/v1"

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
			fmt.Println("Error parsing completion:", ok)
			os.Exit(1)
		}
		fmt.Println(highlight.RegularHighlight(completion))
		s.AddMessage("assistant", completion)
	}
	err = s.Save()
	if err != nil {
		return "", err
	}

	return response, nil
}

func (c *Client) MarshalRequest(prompt string, stream bool, tokens int, model string, system string,
	session *session.Session) (*http.Request, error) {
	messages := prependContext(prompt, system, session)
	payload := CompletionRequest{
		Model:     model,
		Messages:  messages,
		MaxTokens: tokens,
		Stream:    stream,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/chat/completions", v1apiBaseURL), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func prependContext(prompt string, system string, session *session.Session) []Message {
	var messages []Message
	if system != "" {
		messages = append(messages, Message{
			Role:    "system",
			Content: system,
		})
	}
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

func (c *Client) ParseResponse(resp *http.Response) (string, error) {
	var completion CompletionResponse
	err := json.NewDecoder(resp.Body).Decode(&completion)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response payload: %w", err)
	}
	return completion.Choices[0].Message.Content, nil
}

func (c *Client) ParseStreamingResponse(resp *http.Response) (string, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var contentBuilder strings.Builder
	//streamHighlighter := highlight.PartialHighlighter()
	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("error reading response: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if line == "data: [DONE]" {
			break
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			var event CompletionStreamResponse
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				return "", fmt.Errorf("error unmarshaling event: %w", err)
			}
			content := event.Choices[0].Delta.Content
			fmt.Print(content)
			contentBuilder.WriteString(content)
		}
	}

	return contentBuilder.String(), nil
}
