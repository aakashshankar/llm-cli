# LLM CLI 🚀 

A CLI tool to interact with LLMs. Provide your API keys in this format: `<completer>_API_KEY>`

The supported completers are:
- [x] `ANTHROPIC`
- [x] `MISTRAL`
- [ ] `OPENAI`

## Installation ⚙️

Currently, you can build from source and use the generated binary.
Clone the repository and run the following command:

```bash
make
```

This will create a binary named `llm` in the `bin` directory.

## Usage 💻

### 1. CLIs

To use the CLI, you need to provide your API keys as environment variables.
You can set these variables in your shell configuration file (e.g., `.zshrc`) or directly in the terminal.

Once you have set the environment variables, you can use the following command (from repo root) to prompt the LLM:

```bash
bin/llm <llm-name> [flags] <prompt>
```
To chat, simple use the `chat` subcommand:
```bash
bin/llm <llm-name> chat
```
And of course, add `/path/to/repo/bin` to your `$PATH` to use the CLI from anywhere.

### 2. Command assist
Create an environment variable `DEFAULT_COMPLETER` with the name of the completer you want to use. They are:
- `claude`
- `mistral`

Then, you can prefix any command you want to understand with `llm` to understand how it works with examples.