package ui

import (
	"testing"
	// "time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestMainMenuModel_InitialState(t *testing.T) {
	choiceChan := make(chan string, 1)
	defer close(choiceChan) // Ensure channel is closed if test panics or finishes early

	model := NewMainMenuModel(choiceChan)

	expectedItems := []menuItem{
		{title: "Interactive Chat", id: "chat"},
		{title: "Session Management", id: "sessions"},
		{title: "View Configuration", id: "config"},
		{title: "Quit", id: "quit"},
	}

	assert.Equal(t, len(expectedItems), len(model.list.Items()), "Incorrect number of menu items")

	for i, item := range model.list.Items() {
		menuItm, ok := item.(menuItem)
		assert.True(t, ok, "Item is not of type menuItem")
		assert.Equal(t, expectedItems[i].title, menuItm.title, "Menu item title mismatch")
		assert.Equal(t, expectedItems[i].id, menuItm.id, "Menu item ID mismatch")
	}
	assert.Equal(t, "LLM CLI - Main Menu", model.list.Title, "Main menu title mismatch")
}

func TestMainMenuModel_SelectAndEnter(t *testing.T) {
	testCases := []struct {
		name         string
		itemIndex    int
		expectedID   string
	}{
		{name: "Select Interactive Chat", itemIndex: 0, expectedID: "chat"},
		{name: "Select Session Management", itemIndex: 1, expectedID: "sessions"},
		{name: "Select View Configuration", itemIndex: 2, expectedID: "config"},
		{name: "Select Quit", itemIndex: 3, expectedID: "quit"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			choiceChan := make(chan string, 1)
			// No defer close here, as the model closes it upon selection

			model := NewMainMenuModel(choiceChan)
			
			// Select the item
			model.list.Select(tc.itemIndex)
			assert.Equal(t, tc.itemIndex, model.list.Index(), "Item not selected correctly")

			// Send Enter key
			updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = updatedModel.(mainMenuModel)

			assert.True(t, model.quitting, "Model should be quitting after Enter")
			assert.Equal(t, tea.Quit, cmd, "Enter should return tea.Quit command")
			
			// Check choice channel (use a select with timeout to prevent test hanging)
			select {
			case choice := <-choiceChan:
				assert.Equal(t, tc.expectedID, choice, "Incorrect ID sent to choice channel")
			// case <-time.After(100 * time.Millisecond): // Increased timeout for safety
			// 	t.Fatal("Timeout waiting for choice channel")
			default:
				// This case can be hit if the channel is closed before send,
				// or if the send is somehow not immediate. Given the model logic,
				// the channel should have the value. If this fails, there's an issue.
				// Re-check immediately, as sometimes the send might not be instant.
				choice, ok := <-choiceChan
				if !ok || choice != tc.expectedID {
					t.Fatalf("Expected to receive '%s' from choice channel, got '%s' (channel open: %v)", tc.expectedID, choice, ok)
				}

			}
		})
	}
}

func TestMainMenuModel_QuitMessages(t *testing.T) {
	testKeys := []struct {
		name    string
		keyType tea.KeyType
		keyRune rune
	}{
		{name: "Ctrl+C", keyType: tea.KeyCtrlC, keyRune: 0},
		{name: "Esc", keyType: tea.KeyEsc, keyRune: 0},
		// {name: "Q Key", keyType: tea.KeyRunes, keyRune: 'q'}, // 'q' is not a default quit key in this simple model
	}

	for _, tc := range testKeys {
		t.Run(tc.name, func(t *testing.T) {
			choiceChan := make(chan string, 1)
			// No defer close here, model closes it.

			model := NewMainMenuModel(choiceChan)
			
			keyMsg := tea.KeyMsg{Type: tc.keyType}
			if tc.keyType == tea.KeyRunes { // For 'q' if it were added
				keyMsg.Runes = []rune{tc.keyRune}
			}

			updatedModel, cmd := model.Update(keyMsg)
			model = updatedModel.(mainMenuModel)

			assert.True(t, model.quitting, "Model should be quitting after quit key")
			assert.Equal(t, tea.Quit, cmd, "Quit key should return tea.Quit command")

			select {
			case choice := <-choiceChan:
				assert.Equal(t, "quit", choice, "Incorrect ID ('quit') sent to choice channel on quit key")
			// case <-time.After(100 * time.Millisecond):
			// 	t.Fatal("Timeout waiting for 'quit' on choice channel")
			default:
				choice, ok := <-choiceChan
				if !ok || choice != "quit" {
					t.Fatalf("Expected to receive 'quit' from choice channel, got '%s' (channel open: %v)", choice, ok)
				}
			}
		})
	}
}

func TestMainMenuModel_WindowResize(t *testing.T) {
	choiceChan := make(chan string, 1)
	defer close(choiceChan)

	model := NewMainMenuModel(choiceChan)
	
	// Initial list size is based on default width/height of 0,0 from list.New
	// Let's send an initial size to make it concrete.
	model, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = model.(mainMenuModel)


	initialListWidth, initialListHeight := model.list.Size()

	newWidth := 100
	newHeight := 30
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: newWidth, Height: newHeight})
	model = updatedModel.(mainMenuModel)

	currentListWidth, currentListHeight := model.list.Size()

	assert.NotEqual(t, initialListWidth, currentListWidth, "List width should change on resize")
	assert.NotEqual(t, initialListHeight, currentListHeight, "List height should change on resize")
	
	// Expected calculation:
	// msg.Width - appStyle.GetHorizontalPadding()
	// msg.Height - appStyle.GetVerticalPadding() - titleHeight - helpHeight
	// This is complex to assert exactly without duplicating styling logic.
	// Checking for change is a good first step.
	// Let's try a more specific check for width:
	expectedListWidth := newWidth - appStyle.GetHorizontalPadding()
	assert.Equal(t, expectedListWidth, currentListWidth, "List width not correctly updated after resize")
}

func TestMainMenuModel_ChannelClosure(t *testing.T) {
    choiceChan := make(chan string, 1)
    // Do not defer close(choiceChan) here, the model is responsible for closing it.

    model := NewMainMenuModel(choiceChan)
    model.list.Select(0) // Select "Interactive Chat"

    // Simulate Enter key press
    _, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

    // Check if channel is closed and the correct value was sent
    choice, isOpen := <-choiceChan
    assert.False(t, isOpen, "Choice channel should be closed after selection")
    assert.Equal(t, "chat", choice, "Incorrect choice received from closed channel")

	// Test with quit key
	choiceChan2 := make(chan string, 1)
	model2 := NewMainMenuModel(choiceChan2)
	_, _ = model2.Update(tea.KeyMsg{Type: tea.KeyEsc})
	choice2, isOpen2 := <-choiceChan2
	assert.False(t, isOpen2, "Choice channel should be closed after quit key")
	assert.Equal(t, "quit", choice2, "Incorrect choice 'quit' received from closed channel")

}
[end of ui/main_tui_test.go]
