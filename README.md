# LLM CLI ❄️ 

A CLI tool to interact with LLMs. Provide your API keys in this format: `<completer>_API_KEY>`

The supported completers are:
- [x] `ANTHROPIC`
- [ ] `OPENAI`

## Installation ⚙️

Currently, you can build from source and use the generated binary.
Clone the repository and run the following command:

```bash
make
```

This will create a binary named `llm` in the `bin` directory.

## Usage 💻

To use the CLI, you need to provide your API keys as environment variables.
You can set these variables in your shell configuration file (e.g., `.zshrc`) or directly in the terminal.

Once you have set the environment variables, you can use the following command (from repo root) to prompt or chat with the LLM:

```bash
bin/llm <llm-name> <prompt> | chat
```

## Lots of improvements to come! 🚀

- [ ] Add support for other models 
- [x] Persistent context 
- [x] Chat mode
- [ ] Share context across multiple models
- [x] Code cleanup
- [ ] Pretty print the response
- [ ] Support multiple output formats (e.g., JSON, Markdown, HTML)
- [ ] Auto suggest command corrections!
- [ ] Multiple sessions. 
- [ ] Externalize model configs. Automate their maintenance.
- [ ] Stretch goal: Host LLMs on a third party platform.