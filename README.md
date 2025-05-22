# LLM CLI 🚀

A command-line interface for getting instant coding help and technical explanations while maintaining conversational context.

## Supported Models
- ✅ Anthropic (Claude)
- ✅ Mistral AI
- ✅ OpenAI (GPT)

## Features

- **Coding Assistant**: Get concise explanations for programming concepts, debugging help, and code reviews
- **Context-Aware**: Conversations persist locally, enabling follow-up questions and detailed discussions
- **Multi-Model**: Switch between different LLMs while maintaining conversation context
- **Multiple Sessions**: Manage parallel conversations with unique contexts

## TUI Features 🪄

The LLM CLI now features a rich Text User Interface (TUI) for a more interactive and user-friendly experience.

*   **Main Menu**: Running `llm` without any arguments or flags launches the main menu, providing easy navigation to:
    *   Interactive Chat (uses `DEFAULT_COMPLETER`)
    *   Session Management (list and switch sessions)
    *   View Configuration
    *   Quit

*   **Interactive Chat TUI**:
    *   Launched via `llm <completer> chat` or from the main menu.
    *   Provides a full-screen, scrollable interface for your conversations.
    *   Type your message and press `Enter` to send.
    *   Press `Ctrl+C` or `Esc` to exit the chat.

*   **Session Management TUI**:
    *   Launched by `llm session list` or `llm session switch` (when no session UUID is provided).
    *   Displays a list of all your saved sessions, including details like message count, last modification time, and a preview of the first message.
    *   Navigate with arrow keys. If launched via `llm session switch` or from the main menu's "Session Management" option, pressing `Enter` on a session will switch to it.

*   **Configuration Viewing TUI**:
    *   Launched by the new `llm config` command or from the main menu.
    *   Shows the status of your API keys (masked for security), the currently set `DEFAULT_COMPLETER`, and the default model for each LLM variant.
    *   Provides clear instructions on how to set these configurations using environment variables.

## Installation

### Using Go Install
```bash
go install github.com/ahhcash/llm-cli/cmd/llm@latest
```

### Building from Source
```bash
git clone https://github.com/ahhcash/llm-cli
cd llm-cli
make
```

## Quick Start

1. Set API keys:
```bash
export ANTHROPIC_API_KEY="your-key"
export MISTRAL_API_KEY="your-key"
export OPENAI_API_KEY="your-key"
```

2. Set default model:
```bash
export DEFAULT_COMPLETER="claude"  # or "mistral" or "gpt"
```

3. Run `llm` without arguments to open the Main Menu TUI and explore other features like Configuration Viewing.

## Usage

Running `llm` without any arguments or flags launches the **Main Menu TUI**, offering easy navigation to all key features.

### Get Coding Help (Direct Prompting)
For quick, one-off questions or code reviews, you can directly prompt an LLM:
```bash
llm claude "Explain goroutines in Go"
llm openai "Review my function: func add(a, b int) int { return a + b }"
llm mistral "Debug error: fatal: not a git repository"
```
This uses the `assist.Assist` functionality. If you provide a prompt without specifying a completer, it will use your `DEFAULT_COMPLETER`.

### Interactive Chat (TUI)
For an interactive, full-screen chat experience:
```bash
llm claude chat     # Starts an interactive chat TUI with Claude
llm openai chat     # Starts an interactive chat TUI with OpenAI
```
You can also access this feature from the Main Menu. Inside the chat TUI, type your message, press `Enter` to send, and `Ctrl+C` or `Esc` to quit.

### Session Management (TUI)
Manage your conversation sessions:
```bash
llm session list             # Launches a TUI to list all sessions.
llm session switch           # Launches the TUI to select a session to switch to.
llm session switch <uuid>    # Directly switches to the specified session UUID.
llm session inspect <uuid>   # Inspects a session's content.
```
The Session Management TUI allows you to view session details and switch between them easily.

### View Configuration (TUI)
To see your current setup:
```bash
llm config                   # Launches a TUI to display API key status, default completer, and model defaults.
```
This TUI also provides instructions on how to set these configurations.

### Advanced Prompting
Customize your prompts with system messages or specific models:
```bash
llm claude -S "You are a security expert" "Review my SSH config"
llm openai --model gpt-4o "Complex system design"
```

## Troubleshooting
- Check environment variables if API calls fail
- Ensure Go 1.21+ for building
- Verify config directory permissions

## License
MIT License - See LICENSE file