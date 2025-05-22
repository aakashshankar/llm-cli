package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	senderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))  // User messages (Blue)
	botStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))  // LLM messages (Cyan)
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true) // Error messages (Red)
)

type chatModel struct {
	viewport    viewport.Model
	textarea    textarea.Model
	messages    []string // Stores "Sender: message" or "Bot: message"
	spinner     spinner.Model
	isLoading   bool // To control spinner visibility
	err         error

	// For passing information from the command
	initialSystemPrompt string
	llmService          string
}

type (
	errMsg         struct{ err error }
	llmResponseMsg string // To carry LLM responses
)

func (e errMsg) Error() string { return e.err.Error() }

func NewChatModel(initialSystemPrompt, llmService string) chatModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 0 // No limit

	ta.SetWidth(50) // Initial width, will be updated
	ta.SetHeight(1) // Single line input

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle() // No special styling for cursor line
	ta.ShowLineNumbers = false

	vp := viewport.New(50, 5) // Initial size, will be updated
	// vp.YPosition = headerHeight // If we add a header

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))


	// Store initial prompt and service, maybe display them or use later
	var initialMessages []string
	if initialSystemPrompt != "" {
		initialMessages = append(initialMessages, botStyle.Render("System: "+initialSystemPrompt))
	}


	return chatModel{
		textarea:            ta,
		messages:            initialMessages,
		viewport:            vp,
		spinner:             s,
		isLoading:           false,
		err:                 nil,
		initialSystemPrompt: initialSystemPrompt,
		llmService:          llmService,
	}
}

func (m chatModel) Init() tea.Cmd {
	return textarea.Blink // Start the textarea blinking
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		spCmd tea.Cmd
		cmds  []tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	if m.isLoading {
		m.spinner, spCmd = m.spinner.Update(msg)
		cmds = append(cmds, spCmd)
	}


	cmds = append(cmds, tiCmd, vpCmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println("Exiting...")
			return m, tea.Quit
		case tea.KeyEnter:
			userInput := m.textarea.Value()
			if strings.TrimSpace(userInput) == "" {
				return m, nil // Do nothing if input is empty
			}

			m.messages = append(m.messages, senderStyle.Render("You: ")+userInput)
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			
			// Simulate LLM response
			m.isLoading = true
			cmds = append(cmds, m.spinner.Tick) // Start spinner
			// This is where we would make the actual LLM call
			// For now, simulate a delay and a canned response
			cmds = append(cmds, func() tea.Msg {
				time.Sleep(1 * time.Second) // Simulate network latency
				return llmResponseMsg(fmt.Sprintf("LLM (%s): You said: %s", m.llmService, userInput))
			})
			
		}
	case tea.WindowSizeMsg:
		// Example: Header height of 0, Footer (textarea + help) height of say 3-4 lines
		// For now, let's make textarea 1 line and viewport take the rest minus 1 for status/error
		const inputAreaHeight = 1
		const statusAreaHeight = 1
		
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - inputAreaHeight - statusAreaHeight
		m.textarea.SetWidth(msg.Width)
		// m.textarea.SetHeight(inputAreaHeight) // already 1
		m.viewport.SetContent(strings.Join(m.messages, "\n")) // Re-render messages with new width
		m.viewport.GotoBottom()


	case llmResponseMsg:
		m.isLoading = false // Stop spinner
		m.messages = append(m.messages, botStyle.Render(string(msg)))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		
	case errMsg:
		m.err = msg
		m.isLoading = false // Stop spinner if an error occurs
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m chatModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %s\n\n%s\n%s", errorStyle.Render(m.err.Error()), m.viewport.View(), m.textarea.View())
	}

	var loading string
	if m.isLoading {
		loading = m.spinner.View() + " Thinking..."
	}

	// Viewport on top, then loading indicator (if any), then textarea
	return fmt.Sprintf("%s\n%s\n%s", m.viewport.View(), loading, m.textarea.View())
}

// StartChatTUI initializes and runs the chat TUI.
func StartChatTUI(initialSystemPrompt string, llmService string) {
	model := NewChatModel(initialSystemPrompt, llmService)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

[end of ui/chat_tui.go]
