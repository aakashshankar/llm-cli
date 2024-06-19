package gemini

import (
	"context"
	"errors"
	"fmt"
	"github.com/aakashshankar/llm-cli/session"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"net/http"
	"os"
	"strings"
)

// NOTE: Gemini's Go SDK only supports gRPC lol, so the HTTP implementations for response parsing are blank

type Client struct {
	config *Config
	client *genai.Client
}

func NewClient(config *Config) *Client {
	ctx := context.Background()
	geminiClient, err := genai.NewClient(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		fmt.Println("Error creating gemini client:", err)
		os.Exit(1)
	}
	return &Client{
		config: config,
		client: geminiClient,
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

	ctx := context.Background()
	generativeModel := c.client.GenerativeModel(model)
	cs := generativeModel.StartChat()
	cs.History = PrependContext(s)
	promptPart := genai.Text(prompt)
	var contentBuilder strings.Builder
	if stream {
		iter := cs.SendMessageStream(ctx, promptPart)
		for {
			resp, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return "", fmt.Errorf("error sending message: %w", err)
			}
			content := resp.Candidates[0].Content.Parts[0].(genai.Text)
			fmt.Print(content)
			contentBuilder.WriteString(string(content))
		}
		s.AddMessage("assistant", contentBuilder.String())
	} else {
		resp, err := cs.SendMessage(ctx, promptPart)
		if err != nil {
			return "", fmt.Errorf("error sending message: %w", err)
		}
		response := resp.Candidates[0].Content.Parts[0].(genai.Text)
		fmt.Println(response)
		s.AddMessage("assistant", string(response))
	}
	err := s.Save()
	if err != nil {
		return "", err
	}
	return "", nil
}

func (c *Client) ParseResponse(resp *http.Response) (string, error) {
	// NO IMPLEMENTATION
	return "", nil
}

func (c *Client) ParseStreamingResponse(resp *http.Response) (string, error) {
	// NO IMPLEMENTATION
	return "", nil
}

func (c *Client) MarshalRequest(prompt string, stream bool, tokens int, model string, system string, session *session.Session) (*http.Request, error) {
	// NO IMPLEMENTATION
	return nil, nil

}

func PrependContext(session *session.Session) []*genai.Content {
	var parts []*genai.Content
	for _, message := range session.Messages {
		content := genai.Content{
			Role: message.Role,
			Parts: []genai.Part{
				genai.Text(message.Content),
			},
		}
		parts = append(parts, &content)
	}
	return parts
}
