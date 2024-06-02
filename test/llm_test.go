package test

import (
	"github.com/aakashshankar/llm-cli/mocks"
	"github.com/aakashshankar/llm-cli/variants/anthropic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	ass := assert.New(t)

	err := os.Setenv("ANTHROPIC_API_KEY", "test")
	config, err := anthropic.LoadConfig()
	if err != nil {
		t.Error(err)
	}
	ass.Equal(config.APIKey, "test")
}

func TestNewClient(t *testing.T) {
	ass := assert.New(t)

	err := os.Setenv("ANTHROPIC_API_KEY", "test")
	config, err := anthropic.LoadConfig()
	if err != nil {
		t.Error(err)
	}
	client := anthropic.NewClient(config)
	ass.NotNil(client)
}

func TestPrompt(t *testing.T) {
	mockClient := new(mocks.MockClient)
	ass := assert.New(t)
	mockClient.On("Prompt", mock.AnythingOfType("string"), mock.AnythingOfType("bool"), mock.AnythingOfType("int"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("bool")).
		Return("test", nil)
	response, err := mockClient.Prompt("test", false, 1024, "claude-3-sonnet-20240229", "", false)
	if err != nil {
		t.Error(err)
	}
	ass.Equal(response, "test")
}
