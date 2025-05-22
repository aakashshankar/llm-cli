package ui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestChatModel_InitialState(t *testing.T) {
	// Test with no initial system prompt
	model := NewChatModel("", "test-service")
	assert.True(t, model.textarea.Focused(), "Textarea should be focused on init")
	assert.Empty(t, model.messages, "Messages should be empty initially when no system prompt")

	// Test with an initial system prompt
	systemPrompt := "You are a helpful assistant."
	modelWithPrompt := NewChatModel(systemPrompt, "test-service")
	assert.True(t, modelWithPrompt.textarea.Focused())
	assert.Len(t, modelWithPrompt.messages, 1, "Messages should have one system prompt")
	assert.Contains(t, modelWithPrompt.messages[0], systemPrompt, "System prompt message content mismatch")
	assert.Contains(t, modelWithPrompt.messages[0], "System:", "System prompt should be styled as from System")
}

func TestChatModel_SendUserMessage(t *testing.T) {
	model := NewChatModel("", "test-service")
	initialModel, _ := model.Update(nil) // Get initial state
	model = initialModel.(chatModel)

	// Simulate user typing
	model.textarea.SetValue("Hello, world!")
	
	// Simulate pressing Enter
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.NotNil(t, cmd, "Enter key should return a command (for LLM call and spinner)")
	
	model = updatedModel.(chatModel)

	assert.True(t, model.isLoading, "isLoading should be true after sending a message")
	assert.Equal(t, "", model.textarea.Value(), "Textarea should be cleared after sending a message")
	assert.Len(t, model.messages, 1, "Messages should contain one user message")
	assert.Contains(t, model.messages[0], "You: Hello, world!", "User message content mismatch")

	// Simulate receiving LLM response
	// First, handle spinner tick which is part of the batched command
	if cmds, ok := cmd.(tea.BatchMsg); ok {
		for _, c := range cmds {
			if _, ok := c.(spinnerTickCmd); ok { // Assuming spinnerTickCmd is the type for spinner.Tick
				updatedModel, _ = model.Update(c()) // Execute the tick
				model = updatedModel.(chatModel)
				break
			}
		}
	}
	
	// Now, simulate the actual LLM response message
	// We need to find the command that produces llmResponseMsg
	var llmCmd tea.Cmd
	if batchCmd, ok := cmd.(tea.BatchMsg); ok {
		for _, c := range batchCmd {
			// This is tricky because the actual command is an anonymous function.
			// We'll assume the command that results in llmResponseMsg is present.
			// For a more robust test, the LLM call could be an interface.
			// Here, we directly create and send the response message.
			// This also means we can't directly test the time.Sleep in the original command.
		}
	}
	// Since testing the async command is hard, let's directly send the response message
	responseMsg := llmResponseMsg("LLM (test-service): You said: Hello, world!")
	updatedModel, _ = model.Update(responseMsg)
	model = updatedModel.(chatModel)

	assert.False(t, model.isLoading, "isLoading should be false after receiving LLM response")
	assert.Len(t, model.messages, 2, "Messages should contain user message and LLM response")
	assert.Contains(t, model.messages[1], "LLM (test-service): You said: Hello, world!", "LLM response content mismatch")
}

// Helper type for spinner tick command if needed, though spinner.Tick is usually tea.Cmd
type spinnerTickCmd func() tea.Msg


func TestChatModel_QuitMessages(t *testing.T) {
	model := NewChatModel("", "test-service")

	// Test Ctrl+C
	_, cmdCtrlC := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, tea.Quit, cmdCtrlC, "Ctrl+C should return tea.Quit command")
	
	// Test Esc
	_, cmdEsc := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	assert.Equal(t, tea.Quit, cmdEsc, "Esc should return tea.Quit command")
}

func TestChatModel_WindowResize(t *testing.T) {
	model := NewChatModel("", "test-service")
	
	// Initial size might not be set until first WindowSizeMsg
	// Let's send an initial one to set it up
	model, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = model.(chatModel) // type assertion

	newWidth := 100
	newHeight := 30
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: newWidth, Height: newHeight})
	model = updatedModel.(chatModel)

	// Assuming inputAreaHeight = 1, statusAreaHeight = 1 as in chat_tui.go
	expectedViewportHeight := newHeight - 1 - 1 
	
	assert.Equal(t, newWidth, model.viewport.Width, "Viewport width should be updated")
	assert.Equal(t, expectedViewportHeight, model.viewport.Height, "Viewport height should be updated")
	assert.Equal(t, newWidth, model.textarea.Width(), "Textarea width should be updated")
}


func TestChatModel_EmptyInput(t *testing.T) {
	model := NewChatModel("", "test-service")
	initialMessagesCount := len(model.messages)

	model.textarea.SetValue("   ") // Input with only spaces
	
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.Nil(t, cmd, "Enter with empty input should not return a command")
	
	model = updatedModel.(chatModel)
	assert.False(t, model.isLoading, "isLoading should be false with empty input")
	assert.Len(t, model.messages, initialMessagesCount, "Messages count should not change with empty input")
	assert.NotEqual(t, "", model.textarea.Value(), "Textarea should not be cleared with empty input (it retains the space)")
}

// Mock time for LLM response simulation if needed (advanced)
// For now, we are directly sending llmResponseMsg
var originalTimeAfter = time.After
var mockTimeAfter = func(d time.Duration) <-chan time.Time {
	c := make(chan time.Time, 1)
	go func() {
		time.Sleep(1 * time.Millisecond) // Simulate a very short delay for testing
		c <- time.Now()
	}()
	return c
}

// Example of how to use mock time (not fully integrated into SendUserMessage due to complexity of cmd testing)
func TestChatModel_SendUserMessage_WithMockTime(t *testing.T) {
	// timeAfter = mockTimeAfter // Replace standard time.After
	// defer func() { timeAfter = originalTimeAfter }() // Restore original

	// ... test logic ...
	// This is more involved as you'd need to capture the tea.Cmd,
	// check if it's the function returning the llmResponseMsg, and then execute it.
	// The current SendUserMessage test bypasses this by directly sending llmResponseMsg.
}
[end of ui/chat_tui_test.go]
