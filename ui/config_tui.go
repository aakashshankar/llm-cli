package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/ahhcash/llm-cli/defaults" // For GetAllDefaultModels
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	configTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Padding(1, 0, 1, 2)
	configKeyStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("75")) // Light Blue
	configValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	configHelpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Padding(1, 2)
	configSectionStyle = lipgloss.NewStyle().Bold(true).Underline(true).MarginTop(1).MarginBottom(1)
)

// getMaskedEnv retrieves an environment variable and returns a masked version.
func getMaskedEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Not Set") // Red
	}
	if len(val) > 8 { // Mask if long enough
		return fmt.Sprintf("*****%s (Set)", val[len(val)-4:])
	}
	return "Set (short value)" // For very short keys, avoid showing anything
}

type configModel struct {
	viewport viewport.Model
	ready    bool
	content  string
}

func NewConfigModel() configModel {
	// Load configuration data
	var sb strings.Builder

	sb.WriteString(configSectionStyle.Render("API Keys (Environment Variables)") + "\n")
	sb.WriteString(fmt.Sprintf("%s: %s\n", configKeyStyle.Render("ANTHROPIC_API_KEY"), configValueStyle.Render(getMaskedEnv("ANTHROPIC_API_KEY"))))
	sb.WriteString(fmt.Sprintf("%s: %s\n", configKeyStyle.Render("MISTRAL_API_KEY"), configValueStyle.Render(getMaskedEnv("MISTRAL_API_KEY"))))
	sb.WriteString(fmt.Sprintf("%s: %s\n", configKeyStyle.Render("OPENAI_API_KEY"), configValueStyle.Render(getMaskedEnv("OPENAI_API_KEY"))))
	sb.WriteString(configHelpStyle.Render("  To set, use: export KEY_NAME=\"your_api_key\"\n"))

	sb.WriteString(configSectionStyle.Render("Default Completer (Environment Variable)") + "\n")
	defaultCompleter := os.Getenv("DEFAULT_COMPLETER")
	if defaultCompleter == "" {
		defaultCompleter = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Not Set (will use first available variant)")
	} else {
		defaultCompleter = configValueStyle.Render(defaultCompleter)
	}
	sb.WriteString(fmt.Sprintf("%s: %s\n", configKeyStyle.Render("DEFAULT_COMPLETER"), defaultCompleter))
	sb.WriteString(configHelpStyle.Render("  To set, use: export DEFAULT_COMPLETER=\"variant_name\" (e.g., openai, anthropic)\n"))


	sb.WriteString(configSectionStyle.Render("Default Models per Variant (from config file)") + "\n")
	allDefaultModels := defaults.GetAllDefaultModels()
	if len(allDefaultModels) == 0 {
		sb.WriteString("  No default models configured or configuration file not found.\n")
	} else {
		// Sort keys for consistent order if desired (map iteration order is not guaranteed)
		// For now, direct iteration:
		for variant, model := range allDefaultModels {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", configKeyStyle.Render(variant), configValueStyle.Render(model)))
		}
	}
	sb.WriteString(configHelpStyle.Render(fmt.Sprintf("  These are typically managed by the application (e.g. `llm <variant> set-default-model <model_name>`). Config file: %s\n", defaults.GetConfigFilePathForDisplay())))


	sb.WriteString(configSectionStyle.Render("General Instructions") + "\n")
	sb.WriteString("  - API keys and the default completer are read from environment variables.\n")
	sb.WriteString("  - Ensure these variables are set in your shell's configuration file (e.g., .bashrc, .zshrc) for persistence.\n")
	sb.WriteString("  - Changes to environment variables require restarting your shell session or sourcing the config file.\n")


	return configModel{content: sb.String()}
}

func (m configModel) Init() tea.Cmd {
	return nil
}

func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m configModel) headerView() string {
	return configTitleStyle.Render("LLM CLI Configuration")
}

func (m configModel) footerView() string {
	return configHelpStyle.Align(lipgloss.Center).Render("Press 'q' or 'esc' to quit. Use arrow keys to scroll.")
}

func (m configModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

// StartConfigTUI initializes and runs the configuration display TUI.
func StartConfigTUI() {
	model := NewConfigModel()
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running config TUI: %v\n", err)
		os.Exit(1)
	}
}

// The init() function that provided a temporary workaround for 
// defaults.GetConfigFilePathForDisplay has been removed as the function
// is now expected to exist in the defaults package.
[end of ui/config_tui.go]
