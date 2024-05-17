# Claude CLI 🤖

A CLI for Claude's completions. Provide your API keys as the following environment variables:

- `ANTHROPIC_API_KEY`

## Installation ⚙️

Currently, you can build from source and use the generated binary.
Clone the repository and run the following command:

```bash
make build
```

This will create a binary named `claude` in the `bin` directory.

## Usage 💻

To use the CLI, you need to provide your API keys as environment variables.
You can set these variables in your shell configuration file (e.g., `.zshrc`) or directly in the terminal.

Once you have set the environment variables, you can use the `bin/claude ask` command from the repository root to prompt Claude.

## Lots of improvements to come! 💡

- [ ] Add support for other models 
- [ ] Persistent context 
- [ ] Chat mode
- [ ] Share context across multiple models