package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	mainMenuTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Padding(1, 0, 1, 2)
	mainMenuHelpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Padding(1, 0, 0, 2)
)

// menuItem represents an item in the main menu list.
type menuItem struct {
	title string
	id    string // Used to identify the action
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return "" } // No description for main menu items
func (i menuItem) FilterValue() string { return i.title }

// mainMenuDelegate defines how to render items in the list.
type mainMenuDelegate struct{}

func (d mainMenuDelegate) Height() int                               { return 1 } // Title only
func (d mainMenuDelegate) Spacing() int                              { return 0 }
func (d mainMenuDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d mainMenuDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(menuItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	var renderedStr string
	if index == m.Index() { // Currently selected item
		renderedStr = lipgloss.NewStyle().Foreground(lipgloss.Color("202")).Bold(true).Render("> " + str)
	} else { // Other items
		renderedStr = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render("  " + str)
	}
	fmt.Fprint(w, renderedStr)
}

// mainMenuModel is the Bubble Tea model for the main menu.
type mainMenuModel struct {
	list     list.Model
	choice   chan<- string // Channel to send back the user's selection
	quitting bool
}

// NewMainMenuModel initializes a new mainMenuModel.
func NewMainMenuModel(choiceChan chan<- string) mainMenuModel {
	items := []list.Item{
		menuItem{title: "Interactive Chat", id: "chat"},
		menuItem{title: "Session Management", id: "sessions"},
		menuItem{title: "View Configuration", id: "config"},
		menuItem{title: "Quit", id: "quit"},
	}

	l := list.New(items, mainMenuDelegate{}, 0, 0)
	l.Title = "LLM CLI - Main Menu"
	l.Styles.Title = mainMenuTitleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false) // We'll render our own help

	// Customize keys for a simpler menu experience
	l.KeyMap.Quit = key.NewBinding(key.WithKeys("ctrl+c", "esc")) // Will be handled by model directly
	l.KeyMap.Enter = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select"))

	return mainMenuModel{
		list:   l,
		choice: choiceChan,
	}
}

// Init is the first command that will be run.
func (m mainMenuModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model.
func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.quitting {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Adjust list size, accounting for title and help text
		titleHeight := lipgloss.Height(m.list.Title)
		helpHeight := lipgloss.Height(m.helpView())
		m.list.SetSize(msg.Width-appStyle.GetHorizontalPadding(), msg.Height-appStyle.GetVerticalPadding()-titleHeight-helpHeight)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.list.KeyMap.Quit): // q, ctrl+c, esc
			m.quitting = true
			m.choice <- "quit" // Send "quit" choice
			close(m.choice)    // Close the channel
			return m, tea.Quit

		case key.Matches(msg, m.list.KeyMap.Enter):
			selectedItem, ok := m.list.SelectedItem().(menuItem)
			if !ok {
				// Should not happen with current setup
				return m, nil
			}
			m.quitting = true
			m.choice <- selectedItem.id // Send selected item's ID
			close(m.choice)             // Close the channel
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m mainMenuModel) helpView() string {
	return mainMenuHelpStyle.Render("Use arrow keys to navigate. Press Enter to select or Ctrl+C/Esc to quit.")
}

// View renders the UI.
func (m mainMenuModel) View() string {
	if m.quitting {
		return "Exiting...\n"
	}
	
	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		m.list.View(),
		m.helpView(),
	))
}

// StartMainMenuListTUI initializes and runs the main menu TUI.
// It returns the user's selection (e.g., "chat", "sessions", "quit") or an error.
func StartMainMenuListTUI() (string, error) {
	choiceChan := make(chan string, 1) // Buffered channel of size 1

	model := NewMainMenuModel(choiceChan)
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program. This blocks until tea.Quit is received.
	if _, err := p.Run(); err != nil {
		close(choiceChan) // Ensure channel is closed on error
		return "", fmt.Errorf("error running main menu TUI: %w", err)
	}

	// After the program quits, receive the choice from the channel.
	// This will block until a choice is sent or the channel is closed.
	selectedChoice, ok := <-choiceChan
	if !ok {
		// This might happen if the program exited without sending a choice,
		// though current logic ensures a choice or "quit" is always sent.
		return "quit", fmt.Errorf("main menu channel closed without a choice")
	}

	return selectedChoice, nil
}
[end of ui/main_tui.go]
