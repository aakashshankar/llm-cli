package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ahhcash/llm-cli/session" // Import the session package
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle     = lipgloss.NewStyle().Padding(1, 2)
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Padding(0, 0, 1, 0)
	statusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Align(lipgloss.Right) // Green for success
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true) // Red for errors
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Padding(1,0)
)

// sessionItem represents an item in the list.Model
type sessionItem struct {
	session.SessionInfo // Embed SessionInfo
}

// Implement list.Item interface for sessionItem
func (i sessionItem) Title() string {
	star := " "
	if i.IsLatest {
		star = "*"
	}
	return fmt.Sprintf("%s %s (%d msgs)", star, i.UUID, i.MessageCount)
}

func (i sessionItem) Description() string {
	return fmt.Sprintf("Last used: %s. First msg: '%s'", i.LastModified.Format("2006-01-02 15:04"), i.FirstMessage)
}
func (i sessionItem) FilterValue() string { return i.UUID }

// itemDelegate defines how to render items in the list.
type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 2 } // Title and Description
func (d itemDelegate) Spacing() int                              { return 1 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(sessionItem)
	if !ok {
		return
	}

	fn := func(s ...string) string {
		return strings.Join(s, "\n")
	}
	
	var title, desc string
	if index == m.Index() { // Currently selected item
		title = lipgloss.NewStyle().Foreground(lipgloss.Color("202")).Bold(true).Render(i.Title())
		desc = lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Render(i.Description())
	} else { // Other items
		title = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(i.Title())
		desc = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(i.Description())
	}
	
	fmt.Fprint(w, fn(title, desc))
}


// sessionModel is the Bubble Tea model for listing and switching sessions.
type sessionModel struct {
	list         list.Model
	sessions     []session.SessionInfo
	spinner      spinner.Model
	isLoading    bool
	statusMessage string
	errorMessage string
	width        int
	height       int
}

// Custom messages
type sessionsLoadedMsg struct{ items []list.Item }
type switchedSessionMsg struct{ uuid string }
type errorMsg struct{ err error }

func (e errorMsg) Error() string { return e.err.Error() }


// NewSessionModel initializes a new sessionModel.
func NewSessionModel() sessionModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	l := list.New([]list.Item{}, itemDelegate{}, 0, 0)
	l.Title = "Available Chat Sessions"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false) // We'll manage status messages ourselves
	l.SetFilteringEnabled(false) // For now, no filtering
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "switch to session")),
		}
	}


	return sessionModel{
		list:      l,
		spinner:   s,
		isLoading: true,
	}
}

// Init is the first command that will be run.
func (m sessionModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, loadSessionsCmd)
}

// Update handles messages and updates the model.
func (m sessionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-appStyle.GetHorizontalPadding(), msg.Height-appStyle.GetVerticalPadding()-lipgloss.Height(m.list.Title)-4) // Adjust for title and status/help
		return m, nil

	case sessionsLoadedMsg:
		m.isLoading = false
		m.list.SetItems(msg.items)
		m.statusMessage = fmt.Sprintf("%d sessions loaded.", len(msg.items))
		return m, nil

	case switchedSessionMsg:
		m.statusMessage = fmt.Sprintf("Switched to session %s.", msg.uuid)
		return m, tea.Quit // Quit after switching

	case errorMsg:
		m.isLoading = false
		m.errorMessage = msg.Error()
		// Potentially quit or allow retry, for now just display error
		return m, nil

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break // Let the list handle filter input
		}
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c", "esc"))):
			m.statusMessage = "Session selection cancelled."
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if len(m.list.Items()) == 0 || m.isLoading {
				return m, nil
			}
			selectedItem, ok := m.list.SelectedItem().(sessionItem)
			if !ok {
				m.errorMessage = "Error: Could not select session."
				return m, nil
			}
			// Call actual switch logic
			m.isLoading = true // Show spinner while switching (though it's fast)
			m.statusMessage = fmt.Sprintf("Switching to %s...", selectedItem.UUID)
			cmds = append(cmds, m.spinner.Tick)
			cmds = append(cmds, func() tea.Msg {
				session.SwitchSession(selectedItem.UUID) // Assuming this is synchronous and doesn't return error for now
				return switchedSessionMsg{uuid: selectedItem.UUID}
			})
			return m, tea.Batch(cmds...)
		}
	}

	if m.isLoading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI.
func (m sessionModel) View() string {
	if m.errorMessage != "" {
		return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("Error"),
			errorStyle.Render(m.errorMessage),
			helpStyle.Render("Press any key to exit."), // Simplified error exit
		))
	}

	var content string
	if m.isLoading {
		content = lipgloss.JoinVertical(lipgloss.Center, m.spinner.View()+" Loading sessions...", "")
	} else {
		content = m.list.View()
	}
	
	status := ""
	if m.statusMessage != "" {
		status = statusStyle.Width(m.width - appStyle.GetHorizontalPadding()).Render(m.statusMessage)
	}

	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		content,
		status,
		helpStyle.Render(m.list.ShortHelp()),
	))
}

// loadSessionsCmd is a helper command to load sessions.
func loadSessionsCmd() tea.Msg {
	sessions, err := session.ListSessions()
	if err != nil {
		return errorMsg{err: fmt.Errorf("failed to list sessions: %w", err)}
	}

	items := make([]list.Item, len(sessions))
	for i, s := range sessions {
		items[i] = sessionItem{SessionInfo: s}
	}
	return sessionsLoadedMsg{items: items}
}

// StartSessionListTUI initializes and runs the session listing TUI.
func StartSessionListTUI() {
	model := NewSessionModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running session TUI: %v\n", err)
		os.Exit(1)
	}
}
[end of ui/session_tui.go]
