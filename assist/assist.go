package assist

import (
	"fmt"
	"github.com/aakashshankar/llm-cli/defaults"
	"github.com/aakashshankar/llm-cli/llm"
	"os"
)

const system = `You are a terminal assistant. Your input may be just a single command that can be executed on the terminal or a 
question (in the form of natural language) asking about a particular command. 

If you are provided with anything other than these two,you are to respond with:
"Invalid input: <input provided>. Please provide a command or ask me about one."
Do not preface or suffix this response with anything.

Otherwise, you will need to do the following:
1. You will need to provide a concise explanation of the command, its parameters and an example of its usage.
2. Be concise in your answers.

For example,
If the command provided to you as input is "git clone", you would explain that it is used to clone a repository from a remote server and 
provide an example of its usage as "git clone https://github.com/myusername/myrepo.git".

Another example could be, "what does the -c flag do in the git clone command?", you would explain that it is used to specify the 
configuration file to use when cloning the repository.
`

func Assist(prompt string) (string, error) {
	variant, ok := os.LookupEnv("DEFAULT_COMPLETER")
	if !ok {
		fmt.Println("DEFAULT_COMPLETER environment variable not set")
		os.Exit(1)
	}

	defaultModel, err := defaults.GetDefaultModel(variant)
	if err != nil {
		fmt.Println("Error getting default model:", err)
		os.Exit(1)
	}

	newLLM, err := llm.NewLLM(variant)
	if err != nil {
		return "", err
	}

	result, err := newLLM.Prompt(prompt, true, 1024, defaultModel, system, true)
	if err != nil {
		return "", err
	}
	return result, nil
}
