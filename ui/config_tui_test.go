package ui

import (
	"os"
	"strings"
	"testing"

	"github.com/ahhcash/llm-cli/defaults"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// Store original functions to restore them after tests
var originalGetAllDefaultModels func() map[string]string
var originalGetConfigFilePathForDisplay func() string

func setupConfigMocks() {
	originalGetAllDefaultModels = defaults.GetAllDefaultModels
	defaults.GetAllDefaultModels = func() map[string]string {
		return map[string]string{
			"openai":   "gpt-3.5-turbo-test",
			"anthropic": "claude-2-test",
		}
	}

	originalGetConfigFilePathForDisplay = defaults.GetConfigFilePathForDisplay
	defaults.GetConfigFilePathForDisplay = func() string {
		return "/mock/path/to/model_configs.json"
	}
}

func teardownConfigMocks() {
	defaults.GetAllDefaultModels = originalGetAllDefaultModels
	defaults.GetConfigFilePathForDisplay = originalGetConfigFilePathForDisplay
}

func TestConfigModel_InitialContent(t *testing.T) {
	setupConfigMocks()
	defer teardownConfigMocks()

	// Set mock environment variables
	os.Setenv("ANTHROPIC_API_KEY", "anthropic_test_key_123456789")
	os.Setenv("MISTRAL_API_KEY", "") // Test not set
	os.Setenv("OPENAI_API_KEY", "sk-1234") // Test short key
	os.Setenv("DEFAULT_COMPLETER", "openai_test_completer")
	defer func() {
		os.Unsetenv("ANTHROPIC_API_KEY")
		os.Unsetenv("MISTRAL_API_KEY")
		os.Unsetenv("OPENAI_API_KEY")
		os.Unsetenv("DEFAULT_COMPLETER")
	}()

	model := NewConfigModel()

	// Test API Keys
	assert.Contains(t, model.content, "ANTHROPIC_API_KEY: "+configValueStyle.Render("*****6789 (Set)"), "Anthropic key mismatch")
	assert.Contains(t, model.content, "MISTRAL_API_KEY: "+lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Not Set"), "Mistral key mismatch (should be Not Set)")
	assert.Contains(t, model.content, "OPENAI_API_KEY: "+configValueStyle.Render("Set (short value)"), "OpenAI key mismatch (should be short value)")
	
	// Test Default Completer
	assert.Contains(t, model.content, "DEFAULT_COMPLETER: "+configValueStyle.Render("openai_test_completer"), "Default completer mismatch")

	// Test Default Models (from mock)
	assert.Contains(t, model.content, configKeyStyle.Render("openai")+": "+configValueStyle.Render("gpt-3.5-turbo-test"), "OpenAI default model mismatch")
	assert.Contains(t, model.content, configKeyStyle.Render("anthropic")+": "+configValueStyle.Render("claude-2-test"), "Anthropic default model mismatch")
	
	// Test Config File Path (from mock)
	assert.Contains(t, model.content, "/mock/path/to/model_configs.json", "Config file path mismatch")

	// Test Instructions
	assert.Contains(t, model.content, "export KEY_NAME=\"your_api_key\"", "API key instruction missing")
	assert.Contains(t, model.content, "export DEFAULT_COMPLETER=\"variant_name\"", "Default completer instruction missing")
}

func TestConfigModel_QuitMessages(t *testing.T) {
	model := NewConfigModel() // Content doesn't matter for this test

	// Initialize the viewport by sending a WindowSizeMsg, as Update might not process keys otherwise
	_, cmdInit := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	assert.Nil(t, cmdInit, "Initial WindowSizeMsg should not return a command that quits")


	// Test Ctrl+C
	updatedModel, cmdCtrlC := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, tea.Quit, cmdCtrlC, "Ctrl+C should return tea.Quit command")
	model = updatedModel.(configModel)


	// Test Esc
	updatedModel, cmdEsc := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	assert.Equal(t, tea.Quit, cmdEsc, "Esc should return tea.Quit command")
	_ = updatedModel.(configModel)
}

func TestConfigModel_ViewportNavigation(t *testing.T) {
	model := NewConfigModel() // Uses default content which should be long enough

	// Initialize viewport
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 10}) // Small height to ensure scrolling
	model = updatedModel.(configModel)
	assert.True(t, model.ready, "Model should be ready after WindowSizeMsg")

	initialYOffset := model.viewport.YOffset

	// Test KeyDown
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = model.(configModel)
	if model.viewport.AtBottom() {
		assert.Equal(t, initialYOffset, model.viewport.YOffset, "YOffset should not change if already at bottom or not scrollable")
	} else {
		assert.True(t, model.viewport.YOffset > initialYOffset, "YOffset should increase on KeyDown (if not at bottom)")
	}
	
	// Test KeyUp
	currentYOffsetBeforeUp := model.viewport.YOffset
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = model.(configModel)
    if currentYOffsetBeforeUp == 0 { // Was at the top
		assert.Equal(t, 0, model.viewport.YOffset, "YOffset should remain 0 if already at top")
	} else {
		assert.True(t, model.viewport.YOffset < currentYOffsetBeforeUp, "YOffset should decrease on KeyUp (if not at top)")
	}
}

func TestConfigModel_WindowResize(t *testing.T) {
	model := NewConfigModel()

	// Initial size
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(configModel)
	assert.True(t, model.ready, "Model should be ready after initial WindowSizeMsg")
	assert.Equal(t, 80, model.viewport.Width, "Initial viewport width incorrect")
	
	headerHeight := lipgloss.Height(model.headerView())
	footerHeight := lipgloss.Height(model.footerView())
	expectedInitialHeight := 24 - headerHeight - footerHeight
	assert.Equal(t, expectedInitialHeight, model.viewport.Height, "Initial viewport height incorrect")

	// New size
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model = updatedModel.(configModel)
	assert.Equal(t, 100, model.viewport.Width, "Updated viewport width incorrect")
	expectedUpdatedHeight := 30 - headerHeight - footerHeight
	assert.Equal(t, expectedUpdatedHeight, model.viewport.Height, "Updated viewport height incorrect")
}
[end of ui/config_tui_test.go]
