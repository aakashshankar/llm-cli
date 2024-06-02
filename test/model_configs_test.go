package test

import (
	"github.com/aakashshankar/llm-cli/defaults"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultModels(t *testing.T) {
	ass := assert.New(t)
	claude, _ := defaults.GetDefaultModel("claude")
	mistral, _ := defaults.GetDefaultModel("mistral")
	ass.Equal(claude, "claude-3-sonnet-20240229")
	ass.Equal(mistral, "mistral-large-latest")

	_, err := defaults.GetDefaultModel("unknown")
	ass.Error(err)
}
